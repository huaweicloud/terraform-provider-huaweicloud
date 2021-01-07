package huaweicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/security/groups"
)

func resourceIecSecurityGroup() *schema.Resource {

	return &schema.Resource{
		Create: resourceIecSecurityGroupV1Create,
		Read:   resourceIecSecurityGroupV1Read,
		Delete: resourceIecSecurityGroupV1Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ethertype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_range_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remote_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_ip_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIecSecurityGroupV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	group, err := groups.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "HuaweiCloud IEC Security group")
	}

	d.Set("id", group.ID)
	d.Set("name", group.Name)
	d.Set("description", group.Description)

	secRules := make([]map[string]interface{}, len(group.SecurityGroupRules))
	for index, rule := range group.SecurityGroupRules {
		secRules[index] = map[string]interface{}{
			"id":                rule.ID,
			"security_group_id": rule.SecurityGroupID,
			"description":       rule.Description,
			"direction":         rule.Direction,
			"ethertype":         rule.EtherType,
			"protocol":          rule.Protocol,
			"remote_group_id":   rule.RemoteGroupID,
			"remote_ip_prefix":  rule.RemoteIPPrefix,
		}
		if ret, err := strconv.Atoi(rule.PortRangeMax.(string)); err == nil {
			secRules[index]["port_range_max"] = ret
		}
		if ret, err := strconv.Atoi(rule.PortRangeMin.(string)); err == nil {
			secRules[index]["port_range_min"] = ret
		}
	}

	d.Set("security_group_rules", secRules)

	return nil
}

func resourceIecSecurityGroupV1Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	createOpts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	group, err := groups.Create(iecClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC Security Group: %s", err)
	}

	d.SetId(group.ID)
	return resourceIecSecurityGroupV1Read(d, meta)
}

func resourceIecSecurityGroupV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecurityGroupDelete(iecClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud IEC Security Group: %s", err)
	}

	d.SetId("")

	return nil
}

func waitForSecurityGroupDelete(iecClient *golangsdk.ServiceClient, groupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		log.Printf("[DEBUG] Attempting to delete HuaweiCloud Security Group %s.\n", groupID)
		sg, err := groups.Get(iecClient, groupID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud IEC Security Group %s", groupID)
				return sg, "DELETED", nil
			}
			return sg, "ACTIVE", err
		}

		err = groups.Delete(iecClient, groupID).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud IEC Security Group %s", groupID)
				return sg, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return sg, "ACTIVE", nil
				}
			}
			return sg, "ACTIVE", err
		}

		log.Printf("[DEBUG] HuaweiCloud IEC Security Group %s still active.\n", groupID)
		return sg, "ACTIVE", nil
	}
}
