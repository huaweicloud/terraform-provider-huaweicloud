package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/endpointgroups"
)

func resourceVpnEndpointGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpnEndpointGroupV2Create,
		Read:   resourceVpnEndpointGroupV2Read,
		Update: resourceVpnEndpointGroupV2Update,
		Delete: resourceVpnEndpointGroupV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceVpnEndpointGroupV2Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var createOpts endpointgroups.CreateOptsBuilder

	endpointType := resourceVpnEndpointGroupV2EndpointType(d.Get("type").(string))
	v := d.Get("endpoints").([]interface{})
	endpoints := make([]string, len(v))
	for i, v := range v {
		endpoints[i] = v.(string)
	}

	createOpts = VpnEndpointGroupCreateOpts{
		endpointgroups.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			TenantID:    d.Get("tenant_id").(string),
			Endpoints:   endpoints,
			Type:        endpointType,
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create group: %#v", createOpts)

	group, err := endpointgroups.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForEndpointGroupCreation(networkingClient, group.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] EndpointGroup created: %#v", group)

	d.SetId(group.ID)

	return resourceVpnEndpointGroupV2Read(d, meta)
}

func resourceVpnEndpointGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Retrieve information about group: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	group, err := endpointgroups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "group")
	}

	log.Printf("[DEBUG] Read HuaweiCloud Endpoint EndpointGroup %s: %#v", d.Id(), group)

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("tenant_id", group.TenantID)
	d.Set("type", group.Type)
	d.Set("endpoints", group.Endpoints)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceVpnEndpointGroupV2Update(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	opts := endpointgroups.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	var updateOpts endpointgroups.UpdateOptsBuilder
	updateOpts = opts

	log.Printf("[DEBUG] Updating endpoint group with id %s: %#v", d.Id(), updateOpts)

	if hasChange {
		group, err := endpointgroups.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return err
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForEndpointGroupUpdate(networkingClient, group.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForState()

		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Updated group with id %s", d.Id())
	}

	return resourceVpnEndpointGroupV2Read(d, meta)
}

func resourceVpnEndpointGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy group: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	err = endpointgroups.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForEndpointGroupDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()

	return err
}

func waitForEndpointGroupDeletion(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
		log.Printf("[DEBUG] Got group %s => %#v", id, group)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] EndpointGroup %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("Unexpected error: %s", err)
		}

		log.Printf("[DEBUG] EndpointGroup %s deletion is pending", id)
		return group, "DELETING", nil
	}
}

func waitForEndpointGroupCreation(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return group, "ACTIVE", nil
	}
}

func waitForEndpointGroupUpdate(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return group, "UPDATED", nil
	}
}

func resourceVpnEndpointGroupV2EndpointType(epType string) endpointgroups.EndpointType {
	var et endpointgroups.EndpointType
	switch epType {
	case "subnet":
		et = endpointgroups.TypeSubnet
	case "cidr":
		et = endpointgroups.TypeCIDR
	case "vlan":
		et = endpointgroups.TypeVLAN
	case "router":
		et = endpointgroups.TypeRouter
	case "network":
		et = endpointgroups.TypeNetwork
	}
	return et
}
