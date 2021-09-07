package huaweicloud

import (
	"github.com/chnsz/golangsdk/openstack/networking/v1/routetables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceVPCRouteTableV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcRouteTableV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVpcRouteTableV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcRouteTableClient, err := config.RouteTablesV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	logp.Printf("[INFO] Retrieved Vpc Route using given filter %s ", vpcRouteTableClient)

	routetables_id := d.Get("id").(string)

	RouteTables, err := routetables.Get(vpcRouteTableClient, routetables_id).Extract()
	if err != nil {
		return err
	}
	logp.Printf("[DEBUG] Retrieved huaweicloud_vpc_routetable %s: %+v", RouteTables.RouteID, RouteTables)
	d.SetId(RouteTables.RouteID)

	d.Set("name", RouteTables.Name)
	d.Set("destination", RouteTables.Destination)
	d.Set("vpc_id", RouteTables.VPC_ID)
	d.Set("tenant_id", RouteTables.Tenant_Id)
	d.Set("region", GetRegion(d, config))

	routes := make([]map[string]string, 0, len(RouteTables.Routes))
	for _, v := range RouteTables.Routes {
		route := make(map[string]string)
		route["type"] = v.Type
		route["destination"] = v.Destination
		route["nexthop"] = v.Nexthop
		route["system"] = v.System
		routes = append(routes, route)
	}
	d.Set("routes", routes)

	subnets := make([]map[string]string, 0, len(RouteTables.Subnets))
	for _, v := range RouteTables.Subnets {
		subnet := make(map[string]string)
		subnet["id"] = v.Id
		subnets = append(subnets, subnet)
	}
	d.Set("subnets", subnets)

	return nil
}
