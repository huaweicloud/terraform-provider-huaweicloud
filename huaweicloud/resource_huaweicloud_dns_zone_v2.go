package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/dns/v2/zones"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDNSZoneV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSZoneV2Create,
		Read:   resourceDNSZoneV2Read,
		Delete: resourceDNSZoneV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "public",
				ValidateFunc: resourceZoneValidateType,
			},
			"ttl": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      300,
				ValidateFunc: resourceValidateTTL,
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: resourceValidateDescription,
			},
			"masters": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"router": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"router_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"router_region": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"value_specs": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDNSRouter(d *schema.ResourceData) map[string]string {
	router := d.Get("router").(*schema.Set).List()

	if len(router) > 0 {
		mp := make(map[string]string)
		c := router[0].(map[string]interface{})

		if val, ok := c["router_id"]; ok {
			mp["router_id"] = val.(string)
		}
		if val, ok := c["router_region"]; ok {
			mp["router_region"] = val.(string)
		}
		return mp
	}
	return nil
}

func resourceDNSZoneV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	zone_type := d.Get("zone_type").(string)
	router := d.Get("router").(*schema.Set).List()

	// router is required when creating private zone
	if zone_type == "private" {
		if len(router) < 1 {
			return fmt.Errorf("The argument (router) is required when creating HuaweiCloud DNS private zone")
		}
	}
	vs := MapResourceProp(d, "value_specs")
	// Add zone_type to the list.  We do this to keep GopherCloud HuaweiCloud standard.
	vs["zone_type"] = zone_type
	vs["router"] = resourceDNSRouter(d)
	createOpts := ZoneCreateOpts{
		zones.CreateOpts{
			Name:        d.Get("name").(string),
			TTL:         d.Get("ttl").(int),
			Email:       d.Get("email").(string),
			Description: d.Get("description").(string),
		},
		vs,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := zones.Create(dnsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS zone: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS Zone (%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSZone(dnsClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for DNS Zone (%s) to become ACTIVE for creation: %s",
			n.ID, err)
	}

	d.SetId(n.ID)

	log.Printf("[DEBUG] Created HuaweiCloud DNS Zone %s: %#v", n.ID, n)
	return resourceDNSZoneV2Read(d, meta)
}

func resourceDNSZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	n, err := zones.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "zone")
	}

	log.Printf("[DEBUG] Retrieved Zone %s: %#v", d.Id(), n)

	d.Set("name", n.Name)
	d.Set("email", n.Email)
	d.Set("description", n.Description)
	d.Set("ttl", n.TTL)
	if err = d.Set("masters", n.Masters); err != nil {
		return fmt.Errorf("[DEBUG] Error saving masters to state for HuaweiCloud DNS zone (%s): %s", d.Id(), err)
	}
	d.Set("region", GetRegion(d, config))
	d.Set("zone_type", n.ZoneType)

	return nil
}

func resourceDNSZoneV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	var updateOpts zones.UpdateOpts
	if d.HasChange("email") {
		updateOpts.Email = d.Get("email").(string)
	}
	if d.HasChange("ttl") {
		updateOpts.TTL = d.Get("ttl").(int)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	log.Printf("[DEBUG] Updating Zone %s with options: %#v", d.Id(), updateOpts)

	_, err = zones.Update(dnsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud DNS Zone: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS Zone (%s) to update", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSZone(dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	return resourceDNSZoneV2Read(d, meta)
}

func resourceDNSZoneV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	_, err = zones.Delete(dnsClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud DNS Zone: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS Zone (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Target: []string{"DELETED"},
		//we allow to try to delete ERROR zone
		Pending:    []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:    waitForDNSZone(dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for DNS Zone (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceDNSZoneV2ValidType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validTypes := []string{
		"PRIMARY",
		"SECONDARY",
	}

	for _, v := range validTypes {
		if value == v {
			return
		}
	}

	err := fmt.Errorf("%s must be one of %s", k, validTypes)
	errors = append(errors, err)
	return
}

func waitForDNSZone(dnsClient *golangsdk.ServiceClient, zoneId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		zone, err := zones.Get(dnsClient, zoneId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return zone, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud DNS Zone (%s) current status: %s", zone.ID, zone.Status)
		return zone, parseStatus(zone.Status), nil
	}
}

var zoneTypes = [2]string{"public", "private"}

func resourceZoneValidateType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range zoneTypes {
		if value == zoneTypes[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, zoneTypes))

	return
}
