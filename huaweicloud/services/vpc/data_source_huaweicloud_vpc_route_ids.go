package vpc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v2/routes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API VPC GET /v2.0/vpc/routes
func DataSourceVpcRouteIdsV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcRouteIdsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVpcRouteIdsV2Read(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	vpcRouteClient, err := cfg.NetworkingV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating VPC client: %s", err)
	}

	listOpts := routes.ListOpts{
		VPC_ID: d.Get("vpc_id").(string),
	}

	pages, err := routes.List(vpcRouteClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("unable to retrieve VPC routes: %s", err)
	}

	refinedRoutes, err := routes.ExtractRoutes(pages)
	if err != nil {
		return fmt.Errorf("unable to retrieve VPC routes: %s", err)
	}

	if len(refinedRoutes) == 0 {
		return fmt.Errorf("no matching route found for VPC with ID %s", d.Get("vpc_id").(string))
	}

	listRoutes := make([]string, 0)

	for _, route := range refinedRoutes {
		listRoutes = append(listRoutes, route.RouteID)

	}

	d.SetId(d.Get("vpc_id").(string))
	d.Set("ids", listRoutes)
	d.Set("region", cfg.GetRegion(d))

	return nil
}
