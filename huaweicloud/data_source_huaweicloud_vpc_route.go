package huaweicloud

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/routes"

	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceVPCRouteV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcRouteV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nexthop": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVpcRouteV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcRouteClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := routes.ListOpts{
		Type:        d.Get("type").(string),
		Destination: d.Get("destination").(string),
		VPC_ID:      d.Get("vpc_id").(string),
		Tenant_Id:   d.Get("tenant_id").(string),
		RouteID:     d.Get("id").(string),
	}

	pages, err := routes.List(vpcRouteClient, listOpts).AllPages()
	refinedRoutes, err := routes.ExtractRoutes(pages)

	if err != nil {
		return fmt.Errorf("Unable to retrieve vpc routes: %s", err)
	}

	if len(refinedRoutes) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedRoutes) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Route := refinedRoutes[0]

	log.Printf("[INFO] Retrieved Vpc Route using given filter %s: %+v", Route.RouteID, Route)
	d.SetId(Route.RouteID)

	d.Set("type", Route.Type)
	d.Set("nexthop", Route.NextHop)
	d.Set("destination", Route.Destination)
	d.Set("tenant_id", Route.Tenant_Id)
	d.Set("vpc_id", Route.VPC_ID)
	d.Set("id", Route.RouteID)
	d.Set("region", GetRegion(d, config))

	return nil
}
