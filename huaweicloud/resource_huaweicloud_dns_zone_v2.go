package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/dns/v2/zones"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ResourceDNSZoneV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSZoneV2Create,
		Read:   resourceDNSZoneV2Read,
		Update: resourceDNSZoneV2Update,
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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "public",
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: resourceValidateTTL,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: resourceValidateDescription,
			},
			"router": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"router_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"router_region": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"masters": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDNSRouter(d *schema.ResourceData, region string) map[string]string {
	router := d.Get("router").(*schema.Set).List()

	if len(router) > 0 {
		mp := make(map[string]string)
		c := router[0].(map[string]interface{})

		if val, ok := c["router_id"]; ok {
			mp["router_id"] = val.(string)
		}
		if val, ok := c["router_region"]; ok {
			mp["router_region"] = val.(string)
		} else {
			mp["router_region"] = region
		}
		return mp
	}
	return nil
}

func resourceDNSZoneV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	var dnsClient *golangsdk.ServiceClient

	dnsClient, err := config.DnsV2Client(region)
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
		// update the endpoint with region when creating private zone
		dnsClient, err = config.DnsWithRegionClient(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud DNS region client: %s", err)
		}
	}
	vs := MapResourceProp(d, "value_specs")
	// Add zone_type to the list
	vs["zone_type"] = zone_type
	vs["router"] = resourceDNSRouter(d, region)
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

	// router length >1 when creating private zone
	if zone_type == "private" {
		// AssociateZone for the other routers
		routerList := getDNSRouters(d, region)
		if len(routerList) > 1 {
			for i := range routerList {
				// Skip the first router
				if i > 0 {
					log.Printf("[DEBUG] Creating AssociateZone Options: %#v", routerList[i])
					_, err := zones.AssociateZone(dnsClient, n.ID, routerList[i]).Extract()
					if err != nil {
						return fmt.Errorf("Error AssociateZone: %s", err)
					}

					log.Printf("[DEBUG] Waiting for AssociateZone (%s) to Router (%s) become ACTIVE",
						n.ID, routerList[i].RouterID)
					stateRouterConf := &resource.StateChangeConf{
						Target:     []string{"ACTIVE"},
						Pending:    []string{"PENDING"},
						Refresh:    waitForDNSZoneRouter(dnsClient, n.ID, routerList[i].RouterID),
						Timeout:    d.Timeout(schema.TimeoutCreate),
						Delay:      5 * time.Second,
						MinTimeout: 3 * time.Second,
					}

					_, err = stateRouterConf.WaitForState()
					if err != nil {
						return fmt.Errorf("Error waiting for AssociateZone (%s) to Router (%s) become ACTIVE: %s",
							n.ID, routerList[i].RouterID, err)
					}
				} else {
					log.Printf("[DEBUG] First Router Options: %#v", routerList[i])
				}
			}
		}
	}

	d.SetId(n.ID)

	log.Printf("[DEBUG] Created HuaweiCloud DNS Zone %s: %#v", n.ID, n)
	return resourceDNSZoneV2Read(d, meta)
}

func resourceDNSZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	// we can not get the corresponding client by zone type in import scene
	dnsClient, err := config.DnsV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	var zoneInfo *zones.Zone
	zoneInfo, err = zones.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[WARN] fetching zone failed with DNS global endpoint: %s", err)
		// an error occurred while fetching the zone with DNS global endpoint
		// try to fetch it again with DNS region endpoint
		dnsClient, clientErr := config.DnsWithRegionClient(GetRegion(d, config))
		if clientErr != nil {
			// it looks tricky as we return the fetching error rather than clientErr
			return CheckDeleted(d, err, "zone")
		}

		zoneInfo, err = zones.Get(dnsClient, d.Id()).Extract()
		if err != nil {
			return CheckDeleted(d, err, "zone")
		}
	}

	log.Printf("[DEBUG] Retrieved Zone %s: %#v", d.Id(), zoneInfo)

	d.Set("name", zoneInfo.Name)
	d.Set("email", zoneInfo.Email)
	d.Set("description", zoneInfo.Description)
	d.Set("ttl", zoneInfo.TTL)
	if err = d.Set("masters", zoneInfo.Masters); err != nil {
		return fmt.Errorf("[DEBUG] Error saving masters to state for HuaweiCloud DNS zone (%s): %s", d.Id(), err)
	}
	d.Set("region", region)
	d.Set("zone_type", zoneInfo.ZoneType)

	return nil
}

func resourceDNSZoneV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	var dnsClient *golangsdk.ServiceClient

	dnsClient, err := config.DnsV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	zone_type := d.Get("zone_type").(string)
	router := d.Get("router").(*schema.Set).List()

	// router is required when updating private zone
	if zone_type == "private" {
		if len(router) < 1 {
			return fmt.Errorf("The argument (router) is required when updating HuaweiCloud DNS private zone")
		}
		// update the endpoint with region when creating private zone
		dnsClient, err = config.DnsWithRegionClient(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud DNS region client: %s", err)
		}
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

	if d.HasChange("router") {
		// when updating private zone
		if zone_type == "private" {
			associateList, disassociateList, err := resourceGetDNSRouters(dnsClient, d, region)
			if err != nil {
				return fmt.Errorf("Error getting HuaweiCloud DNS Zone Router: %s", err)
			}
			if len(associateList) > 0 {
				// AssociateZone
				for i := range associateList {
					log.Printf("[DEBUG] Updating AssociateZone Options: %#v", associateList[i])
					_, err := zones.AssociateZone(dnsClient, d.Id(), associateList[i]).Extract()
					if err != nil {
						return fmt.Errorf("Error AssociateZone: %s", err)
					}

					log.Printf("[DEBUG] Waiting for AssociateZone (%s) to Router (%s) become ACTIVE",
						d.Id(), associateList[i].RouterID)
					stateRouterConf := &resource.StateChangeConf{
						Target:     []string{"ACTIVE"},
						Pending:    []string{"PENDING"},
						Refresh:    waitForDNSZoneRouter(dnsClient, d.Id(), associateList[i].RouterID),
						Timeout:    d.Timeout(schema.TimeoutUpdate),
						Delay:      5 * time.Second,
						MinTimeout: 3 * time.Second,
					}

					_, err = stateRouterConf.WaitForState()
					if err != nil {
						return fmt.Errorf("Error waiting for AssociateZone (%s) to Router (%s) become ACTIVE: %s",
							d.Id(), associateList[i].RouterID, err)
					}
				}
			}
			if len(disassociateList) > 0 {
				// DisassociateZone
				for j := range disassociateList {
					log.Printf("[DEBUG] Updating DisassociateZone Options: %#v", disassociateList[j])
					_, err := zones.DisassociateZone(dnsClient, d.Id(), disassociateList[j]).Extract()
					if err != nil {
						return fmt.Errorf("Error DisassociateZone: %s", err)
					}

					log.Printf("[DEBUG] Waiting for DisassociateZone (%s) to Router (%s) become DELETED",
						d.Id(), disassociateList[j].RouterID)
					stateRouterConf := &resource.StateChangeConf{
						Target:     []string{"DELETED"},
						Pending:    []string{"ACTIVE", "PENDING", "ERROR"},
						Refresh:    waitForDNSZoneRouter(dnsClient, d.Id(), disassociateList[j].RouterID),
						Timeout:    d.Timeout(schema.TimeoutUpdate),
						Delay:      5 * time.Second,
						MinTimeout: 3 * time.Second,
					}

					_, err = stateRouterConf.WaitForState()
					if err != nil {
						return fmt.Errorf("Error waiting for DisassociateZone (%s) to Router (%s) become DELETED: %s",
							d.Id(), disassociateList[j].RouterID, err)
					}
				}
			}
		}
	}

	return resourceDNSZoneV2Read(d, meta)
}

func resourceDNSZoneV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	var dnsClient *golangsdk.ServiceClient
	var err error

	zoneType := d.Get("zone_type").(string)
	// update the endpoint with region when creating private zone
	if zoneType == "private" {
		dnsClient, err = config.DnsWithRegionClient(GetRegion(d, config))
	} else {
		dnsClient, err = config.DnsV2Client(GetRegion(d, config))
	}
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

func getDNSRouters(d *schema.ResourceData, region string) []zones.RouterOpts {
	router := d.Get("router").(*schema.Set).List()
	if len(router) == 0 {
		return nil
	}

	res := make([]zones.RouterOpts, len(router))
	for i := range router {
		ro := zones.RouterOpts{}
		c := router[i].(map[string]interface{})
		if val, ok := c["router_id"]; ok {
			ro.RouterID = val.(string)
		}
		if val, ok := c["router_region"]; ok {
			ro.RouterRegion = val.(string)
		} else {
			ro.RouterRegion = region
		}

		res[i] = ro
	}
	return res
}

func waitForDNSZoneRouter(dnsClient *golangsdk.ServiceClient, zoneId string, routerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		zone, err := zones.Get(dnsClient, zoneId).Extract()
		if err != nil {
			return nil, "", err
		}
		for i := range zone.Routers {
			if routerId == zone.Routers[i].RouterID {
				log.Printf("[DEBUG] HuaweiCloud DNS Zone (%s) Router (%s) current status: %s",
					zoneId, routerId, zone.Routers[i].Status)
				return zone, parseStatus(zone.Routers[i].Status), nil
			}
		}
		return zone, "DELETED", nil
	}
}

func resourceGetDNSRouters(dnsClient *golangsdk.ServiceClient, d *schema.ResourceData,
	region string) ([]zones.RouterOpts, []zones.RouterOpts, error) {

	// get zone info from api
	n, err := zones.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		return nil, nil, CheckDeleted(d, err, "zone")
	}
	// get routers from local
	localRouters := getDNSRouters(d, region)

	// get associateMap
	associateMap := make(map[string]zones.RouterOpts)
	for _, local := range localRouters {
		// Check if local is found in api
		found := false
		for _, raw := range n.Routers {
			if local.RouterID == raw.RouterID {
				found = true
				break
			}
		}
		// If local is not found in api
		if !found {
			associateMap[local.RouterID] = local
		}
	}

	// convert associateMap to associateList
	associateList := make([]zones.RouterOpts, len(associateMap))
	var i = 0
	for _, associateRouter := range associateMap {
		associateList[i] = associateRouter
		i++
	}

	// get disassociateMap
	disassociateMap := make(map[string]zones.RouterOpts)
	for _, raw := range n.Routers {
		// Check if api is found in local
		found := false
		for _, local := range localRouters {
			if raw.RouterID == local.RouterID {
				found = true
				break
			}
		}
		// If api is not found in local
		if !found {
			disassociateMap[raw.RouterID] = zones.RouterOpts{
				RouterID:     raw.RouterID,
				RouterRegion: raw.RouterRegion,
			}
		}
	}

	// convert disassociateMap to disassociateList
	disassociateList := make([]zones.RouterOpts, len(disassociateMap))
	var j = 0
	for _, disassociateRouter := range disassociateMap {
		disassociateList[j] = disassociateRouter
		j++
	}

	return associateList, disassociateList, nil
}
