package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/firewalls"
)

func resourceIecNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceIecNetworkACLCreate,
		Read:   resourceIecNetworkACLRead,
		Update: resourceIecNetworkACLUpdate,
		Delete: resourceIecNetworkACLDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"networks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"inbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"outbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIecNetworkACLNetworks(d *schema.ResourceData) []firewalls.ReqSubnet {
	rawNetworks := d.Get("networks").(*schema.Set).List()
	networkOpts := make([]firewalls.ReqSubnet, len(rawNetworks))

	for i, val := range rawNetworks {
		raw := val.(map[string]interface{})
		networkOpts[i] = firewalls.ReqSubnet{
			ID:    raw["subnet_id"].(string),
			VpcID: raw["vpc_id"].(string),
		}
	}
	return networkOpts
}

func resourceIecNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	createOpts := firewalls.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	log.Printf("[DEBUG] Create IEC network acl: %#v", createOpts)
	group, err := firewalls.Create(iecClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating IEC network acl: %s", err)
	}
	d.SetId(group.ID)

	// associate subnets
	subnetsOpts := resourceIecNetworkACLNetworks(d)
	if len(subnetsOpts) > 0 {
		updateOpts := firewalls.UpdateOpts{
			Name:    d.Get("name").(string),
			Subnets: &subnetsOpts,
		}
		log.Printf("[DEBUG] attempt to associate IEC network acl with subnets: %#v", updateOpts)
		_, err := firewalls.Update(iecClient, group.ID, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error associating subnets with IEC network acl: %s", err)
		}
	}

	return resourceIecNetworkACLRead(d, meta)
}

func resourceIecNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	fwGroup, err := firewalls.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "iec network acl")
	}

	log.Printf("[DEBUG] Read HuaweiCloud IEC network acl %s: %#v", d.Id(), fwGroup)
	d.Set("name", fwGroup.Name)
	d.Set("status", fwGroup.Status)
	d.Set("description", fwGroup.Description)
	d.Set("inbound_rules", getFirewallRuleIDs(fwGroup.IngressFWPolicy))
	d.Set("outbound_rules", getFirewallRuleIDs(fwGroup.EgressFWPolicy))
	var networkSet []map[string]interface{}
	for _, val := range fwGroup.Subnets {
		subnet := make(map[string]interface{})
		subnet["vpc_id"] = val.VpcID
		subnet["subnet_id"] = val.ID
		networkSet = append(networkSet, subnet)
	}
	if err = d.Set("networks", networkSet); err != nil {
		return fmt.Errorf("Saving iec networks failed: %s", err)
	}

	return nil
}

func resourceIecNetworkACLUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	if d.HasChanges("name", "description", "networks") {
		opts := firewalls.UpdateOpts{
			Name: d.Get("name").(string),
		}

		if d.HasChange("description") {
			desc := d.Get("description").(string)
			opts.Description = &desc
		}

		if d.HasChange("networks") {
			subnetsOpts := resourceIecNetworkACLNetworks(d)
			opts.Subnets = &subnetsOpts
		}

		log.Printf("[DEBUG] Updating IEC network acl with id %s: %#v", d.Id(), opts)
		_, err := firewalls.Update(iecClient, d.Id(), opts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating IEC network acl: %s", err)
		}
	}

	return resourceIecNetworkACLRead(d, meta)
}

func resourceIecNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	// unbind subents before deleting the network acl
	rawNetworks := d.Get("networks").(*schema.Set).List()
	if len(rawNetworks) > 0 {
		unbindSubnets := make([]firewalls.ReqSubnet, 0)
		opts := firewalls.UpdateOpts{
			Name:    d.Get("name").(string),
			Subnets: &unbindSubnets,
		}
		err = firewalls.Update(iecClient, d.Id(), opts).Err
		if err != nil {
			return fmt.Errorf("Error disassociating all subents with IEC network acl: %s", err)
		}
	}

	err = firewalls.Delete(iecClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud iec firewall: %s", err)
	}

	d.SetId("")
	return nil
}
