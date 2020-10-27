package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/security/securitygroups"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/rules"
)

func ResourceNetworkingSecGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingSecGroupV2Create,
		Read:   resourceNetworkingSecGroupV2Read,
		Update: resourceNetworkingSecGroupV2Update,
		Delete: resourceNetworkingSecGroupV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"delete_default_rules": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingSecGroupV2Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	networkingClient, err := config.SecurityGroupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	opts := securitygroups.CreateOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
	}

	log.Printf("[DEBUG] Create HuaweiCloud Neutron Security Group: %#v", opts)

	security_group, err := securitygroups.Create(networkingClient, opts).Extract()
	if err != nil {
		return err
	}

	// Delete the default security group rules if it has been requested.
	deleteDefaultRules := d.Get("delete_default_rules").(bool)
	if deleteDefaultRules {
		networkingClient_del, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		security_group, err := groups.Get(networkingClient_del, security_group.ID).Extract()
		if err != nil {
			return err
		}
		for _, rule := range security_group.Rules {
			if err := rules.Delete(networkingClient_del, rule.ID).ExtractErr(); err != nil {
				return fmt.Errorf(
					"There was a problem deleting a default security group rule: %s", err)
			}
		}
	}

	log.Printf("[DEBUG] HuaweiCloud Neutron Security Group created: %#v", security_group)

	d.SetId(security_group.ID)

	description := d.Get("description").(string)
	networkingClient_des, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}
	updateOpts := groups.UpdateOpts{
		Description: &description,
	}

	_, err = groups.Update(networkingClient_des, d.Id(), updateOpts).Extract()

	return resourceNetworkingSecGroupV2Read(d, meta)
}

func resourceNetworkingSecGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Retrieve information about security group: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.SecurityGroupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	security_group, err := securitygroups.Get(networkingClient, d.Id()).Extract()

	if err != nil {
		return CheckDeleted(d, err, "HuaweiCloud Neutron Security group")
	}

	d.Set("description", security_group.Description)
	d.Set("name", security_group.Name)
	d.Set("region", GetRegion(d, config))
	d.Set("enterprise_project_id", security_group.EnterpriseProjectId)

	return nil
}

func resourceNetworkingSecGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var update bool
	var updateOpts groups.UpdateOpts

	if d.HasChange("name") {
		update = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		update = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if update {
		log.Printf("[DEBUG] Updating SecGroup %s with options: %#v", d.Id(), updateOpts)
		_, err = groups.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud SecGroup: %s", err)
		}
	}

	return resourceNetworkingSecGroupV2Read(d, meta)
}

func resourceNetworkingSecGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy security group: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.SecurityGroupV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecGroupDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud Neutron Security Group: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSecGroupDelete(networkingClient *golangsdk.ServiceClient, secGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud Security Group %s.\n", secGroupId)

		r, err := securitygroups.Get(networkingClient, secGroupId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud Neutron Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = securitygroups.Delete(networkingClient, secGroupId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud Neutron Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] HuaweiCloud Neutron Security Group %s still active.\n", secGroupId)
		return r, "ACTIVE", nil
	}
}
