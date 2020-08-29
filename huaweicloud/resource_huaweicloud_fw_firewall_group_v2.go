package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/routerinsertion"
)

func resourceFWFirewallGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWFirewallGroupV2Create,
		Read:   resourceFWFirewallGroupV2Read,
		Update: resourceFWFirewallGroupV2Update,
		Delete: resourceFWFirewallGroupV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_network_acl resource instead",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ingress_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"egress_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceFWFirewallGroupV2Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	fwClient, err := config.fwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	var createOpts firewall_groups.CreateOptsBuilder

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts = FirewallGroupCreateOpts{
		firewall_groups.CreateOpts{
			Name:            d.Get("name").(string),
			Description:     d.Get("description").(string),
			IngressPolicyID: d.Get("ingress_policy_id").(string),
			EgressPolicyID:  d.Get("egress_policy_id").(string),
			AdminStateUp:    &adminStateUp,
			TenantID:        d.Get("tenant_id").(string),
		},
		MapValueSpecs(d),
	}

	portsRaw := d.Get("ports").(*schema.Set).List()
	if len(portsRaw) > 0 {
		log.Printf("[DEBUG] Will attempt to associate Firewall group with port(s): %+v", portsRaw)

		var portIds []string
		for _, v := range portsRaw {
			portIds = append(portIds, v.(string))
		}

		createOpts = &routerinsertion.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			PortIDs:           portIds,
		}
	}

	log.Printf("[DEBUG] Create firewall group: %#v", createOpts)

	firewall_group, err := firewall_groups.Create(fwClient, createOpts).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Firewall group created: %#v", firewall_group)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForFirewallGroupActive(fwClient, firewall_group.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()
	log.Printf("[DEBUG] Firewall group (%s) is active.", firewall_group.ID)

	d.SetId(firewall_group.ID)

	return resourceFWFirewallGroupV2Read(d, meta)
}

func resourceFWFirewallGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Retrieve information about firewall: %s", d.Id())

	config := meta.(*Config)
	fwClient, err := config.fwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	var firewall_group FirewallGroup
	err = firewall_groups.Get(fwClient, d.Id()).ExtractInto(&firewall_group)
	if err != nil {
		return CheckDeleted(d, err, "firewall")
	}

	log.Printf("[DEBUG] Read HuaweiCloud Firewall group %s: %#v", d.Id(), firewall_group)

	d.Set("name", firewall_group.Name)
	d.Set("description", firewall_group.Description)
	d.Set("ingress_policy_id", firewall_group.IngressPolicyID)
	d.Set("egress_policy_id", firewall_group.EgressPolicyID)
	d.Set("admin_state_up", firewall_group.AdminStateUp)
	d.Set("tenant_id", firewall_group.TenantID)
	if err := d.Set("ports", firewall_group.PortIDs); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ports to state for HuaweiCloud firewall group (%s): %s", d.Id(), err)
	}
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceFWFirewallGroupV2Update(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	fwClient, err := config.fwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	// PolicyID is required
	opts := firewall_groups.UpdateOpts{
		IngressPolicyID: d.Get("ingress_policy_id").(string),
		EgressPolicyID:  d.Get("egress_policy_id").(string),
	}

	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}

	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &adminStateUp
	}

	var updateOpts firewall_groups.UpdateOptsBuilder
	var portIds []string
	if d.HasChange("ports") {
		portsRaw := d.Get("ports").(*schema.Set).List()
		log.Printf("[DEBUG] Will attempt to associate Firewall group with port(s): %+v", portsRaw)
		for _, v := range portsRaw {
			portIds = append(portIds, v.(string))
		}

		updateOpts = routerinsertion.UpdateOptsExt{
			UpdateOptsBuilder: opts,
			PortIDs:           portIds,
		}
	} else {
		updateOpts = opts
	}

	log.Printf("[DEBUG] Updating firewall with id %s: %#v", d.Id(), updateOpts)

	err = firewall_groups.Update(fwClient, d.Id(), updateOpts).Err
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForFirewallGroupActive(fwClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()

	return resourceFWFirewallGroupV2Read(d, meta)
}

func resourceFWFirewallGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy firewall group: %s", d.Id())

	config := meta.(*Config)
	fwClient, err := config.fwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	// Ensure the firewall group was fully created/updated before being deleted.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForFirewallGroupActive(fwClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()

	err = firewall_groups.Delete(fwClient, d.Id()).Err

	if err != nil {
		return err
	}

	stateConf = &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForFirewallGroupDeletion(fwClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()

	return err
}

func waitForFirewallGroupActive(fwClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		var fw FirewallGroup

		err := firewall_groups.Get(fwClient, id).ExtractInto(&fw)
		if err != nil {
			return nil, "", err
		}
		return fw, fw.Status, nil
	}
}

func waitForFirewallGroupDeletion(fwClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		fw, err := firewall_groups.Get(fwClient, id).Extract()
		log.Printf("[DEBUG] Got firewall group %s => %#v", id, fw)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Firewall group %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("Unexpected error: %s", err)
		}

		log.Printf("[DEBUG] Firewall group %s deletion is pending", id)
		return fw, "DELETING", nil
	}
}
