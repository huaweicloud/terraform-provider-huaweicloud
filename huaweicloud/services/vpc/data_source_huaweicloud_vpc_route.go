package vpc

import (
	"github.com/chnsz/golangsdk/openstack/networking/v2/routes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

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
			"nexthop": {
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
		},
	}
}

func dataSourceVpcRouteV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Vpc client: %s", err)
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
		return fmtp.Errorf("Unable to retrieve vpc routes: %s", err)
	}

	refinedRoutes, err := routes.ExtractRoutes(pages)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve vpc routes: %s", err)
	}

	if len(refinedRoutes) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedRoutes) > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Route := refinedRoutes[0]

	logp.Printf("[INFO] Retrieved Vpc Route using given filter %s: %+v", Route.RouteID, Route)
	d.SetId(Route.RouteID)

	d.Set("type", Route.Type)
	d.Set("nexthop", Route.NextHop)
	d.Set("destination", Route.Destination)
	d.Set("tenant_id", Route.Tenant_Id)
	d.Set("vpc_id", Route.VPC_ID)
	d.Set("id", Route.RouteID)
	d.Set("region", config.GetRegion(d))

	return nil
}
