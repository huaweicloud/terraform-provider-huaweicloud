package huaweicloud

import (
	"github.com/chnsz/golangsdk/openstack/networking/v1/routetables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceVPCRouteTableListV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcRouteTableListV2Read,

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
				//Elem:     &schema.Schema{Type: schema.TypeString},
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
				//Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceVpcRouteTableListV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcRouteTableListClient, err := config.RouteTablesV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := routetables.ListOpts{}

	if v, ok := d.GetOk("vpc_id"); ok {
		listOpts.VPC_ID = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		listOpts.SubnetID = v.(string)
	}

	pages, err := routetables.List(vpcRouteTableListClient, listOpts).AllPages()

	RouteTableLists, err := routetables.ExtractRouteTables(pages)
	if err != nil {
		return fmtp.Errorf("Unable to list huaweicloud_vpc_routetable_list: %s", err)
	}

	RouteTableList := RouteTableLists[0]
	logp.Printf("[DEBUG] Retrieved huaweicloud_vpc_routetable_list %s: %+v", RouteTableList.RouteID, RouteTableList)
	d.SetId(RouteTableList.RouteID)

	d.Set("name", RouteTableList.Name)
	d.Set("destination", RouteTableList.Destination)
	d.Set("vpc_id", RouteTableList.VPC_ID)
	d.Set("tenant_id", RouteTableList.Tenant_Id)
	d.Set("region", GetRegion(d, config))

	return nil
}
