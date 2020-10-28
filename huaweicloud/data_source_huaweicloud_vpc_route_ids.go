package huaweicloud

import (
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/routes"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVPCRouteIdsV2() *schema.Resource {
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
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceVpcRouteIdsV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcRouteClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := routes.ListOpts{
		VPC_ID: d.Get("vpc_id").(string),
	}

	pages, err := routes.List(vpcRouteClient, listOpts).AllPages()
	refinedRoutes, err := routes.ExtractRoutes(pages)

	if err != nil {
		return fmt.Errorf("Unable to retrieve vpc Routes: %s", err)
	}

	if len(refinedRoutes) == 0 {
		return fmt.Errorf("no matching route found for vpc with id %s", d.Get("vpc_id").(string))
	}

	listRoutes := make([]string, 0)

	for _, route := range refinedRoutes {
		listRoutes = append(listRoutes, route.RouteID)

	}

	d.SetId(d.Get("vpc_id").(string))
	d.Set("ids", listRoutes)
	d.Set("region", GetRegion(d, config))

	return nil
}
