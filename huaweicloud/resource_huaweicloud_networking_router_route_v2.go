package huaweicloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceNetworkingRouterRouteV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingRouterRouteV2Create,
		Read:   resourceNetworkingRouterRouteV2Read,
		Delete: resourceNetworkingRouterRouteV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_vpc_route resource instead",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingRouterRouteV2Create(d *schema.ResourceData, meta interface{}) error {

	routerId := d.Get("router_id").(string)
	config.MutexKV.Lock(routerId)
	defer config.MutexKV.Unlock(routerId)

	var destCidr string = d.Get("destination_cidr").(string)
	var nextHop string = d.Get("next_hop").(string)

	conf := meta.(*config.Config)
	networkingClient, err := conf.NetworkingV2Client(GetRegion(d, conf))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	n, err := routers.Get(networkingClient, routerId).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud Neutron Router: %s", err)
	}

	var updateOpts routers.UpdateOpts
	var routeExists bool = false

	var rts []routers.Route = n.Routes
	for _, r := range rts {

		if r.DestinationCIDR == destCidr && r.NextHop == nextHop {
			routeExists = true
			break
		}
	}

	if !routeExists {

		if destCidr != "" && nextHop != "" {
			r := routers.Route{DestinationCIDR: destCidr, NextHop: nextHop}
			logp.Printf(
				"[INFO] Adding route %s", r)
			rts = append(rts, r)
		}

		updateOpts.Routes = rts

		logp.Printf("[DEBUG] Updating Router %s with options: %+v", routerId, updateOpts)

		_, err = routers.Update(networkingClient, routerId, updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud Neutron Router: %s", err)
		}
		d.SetId(fmt.Sprintf("%s-route-%s-%s", routerId, destCidr, nextHop))

	} else {
		logp.Printf("[DEBUG] Router %s has route already", routerId)
	}

	return resourceNetworkingRouterRouteV2Read(d, meta)
}

func resourceNetworkingRouterRouteV2Read(d *schema.ResourceData, meta interface{}) error {

	routerId := d.Get("router_id").(string)

	conf := meta.(*config.Config)
	networkingClient, err := conf.NetworkingV2Client(GetRegion(d, conf))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	destCidr := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)

	if d.Id() != "" && strings.Contains(d.Id(), "-route-") {
		routeIDParts := strings.Split(d.Id(), "-route-")
		routeLastIDParts := strings.Split(routeIDParts[1], "-")

		if routerId == "" {
			routerId = routeIDParts[0]
			d.Set("router_id", routerId)
		}
		if destCidr == "" {
			destCidr = routeLastIDParts[0]
		}
		if nextHop == "" {
			nextHop = routeLastIDParts[1]
		}
	}

	n, err := routers.Get(networkingClient, routerId).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud Neutron Router: %s", err)
	}

	logp.Printf("[DEBUG] Retrieved Router %s: %+v", routerId, n)

	d.Set("next_hop", "")
	d.Set("destination_cidr", "")

	for _, r := range n.Routes {

		if r.DestinationCIDR == destCidr && r.NextHop == nextHop {
			d.Set("destination_cidr", destCidr)
			d.Set("next_hop", nextHop)
			break
		}
	}

	d.Set("region", GetRegion(d, conf))

	return nil
}

func resourceNetworkingRouterRouteV2Delete(d *schema.ResourceData, meta interface{}) error {

	routerId := d.Get("router_id").(string)
	config.MutexKV.Lock(routerId)
	defer config.MutexKV.Unlock(routerId)

	conf := meta.(*config.Config)

	networkingClient, err := conf.NetworkingV2Client(GetRegion(d, conf))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	n, err := routers.Get(networkingClient, routerId).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud Neutron Router: %s", err)
	}

	var updateOpts routers.UpdateOpts

	var destCidr string = d.Get("destination_cidr").(string)
	var nextHop string = d.Get("next_hop").(string)

	var oldRts []routers.Route = n.Routes
	var newRts []routers.Route

	for _, r := range oldRts {

		if r.DestinationCIDR != destCidr || r.NextHop != nextHop {
			newRts = append(newRts, r)
		}
	}

	if len(oldRts) != len(newRts) {
		r := routers.Route{DestinationCIDR: destCidr, NextHop: nextHop}
		logp.Printf(
			"[INFO] Deleting route %s", r)
		updateOpts.Routes = newRts

		logp.Printf("[DEBUG] Updating Router %s with options: %+v", routerId, updateOpts)

		_, err = routers.Update(networkingClient, routerId, updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud Neutron Router: %s", err)
		}
	} else {
		return fmtp.Errorf("Route did not exist already")
	}

	return nil
}
