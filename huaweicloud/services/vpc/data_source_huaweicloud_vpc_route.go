package vpc

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v2/routes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API VPC GET /v2.0/vpc/routes
func DataSourceVpcRouteV2() *schema.Resource {
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
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nexthop": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcRouteV2Read(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	vpcRouteClient, err := cfg.NetworkingV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating VPC client: %s", err)
	}

	listOpts := routes.ListOpts{
		Type:        d.Get("type").(string),
		Destination: d.Get("destination").(string),
		VPC_ID:      d.Get("vpc_id").(string),
		Tenant_Id:   d.Get("tenant_id").(string),
		RouteID:     d.Get("id").(string),
	}

	pages, err := routes.List(vpcRouteClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("unable to retrieve VPC routes: %s", err)
	}

	refinedRoutes, err := routes.ExtractRoutes(pages)
	if err != nil {
		return fmt.Errorf("unable to retrieve VPC routes: %s", err)
	}

	if len(refinedRoutes) < 1 {
		return errors.New("your query returned no results, please change your search criteria and try again")
	}

	if len(refinedRoutes) > 1 {
		return errors.New("your query returned more than one result, please try a more specific search criteria")
	}

	Route := refinedRoutes[0]

	log.Printf("[INFO] Retrieved VPC route using given filter %s: %+v", Route.RouteID, Route)
	d.SetId(Route.RouteID)

	d.Set("type", Route.Type)
	d.Set("nexthop", Route.NextHop)
	d.Set("destination", Route.Destination)
	d.Set("tenant_id", Route.Tenant_Id)
	d.Set("vpc_id", Route.VPC_ID)
	d.Set("id", Route.RouteID)
	d.Set("region", cfg.GetRegion(d))

	return nil
}
