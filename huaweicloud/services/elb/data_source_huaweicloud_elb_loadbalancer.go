package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceElbLoadbalancesV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbLoadbalancersV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"share_type": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"dedicated", "share",
				}, false),
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv4_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"l4_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"l7_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loadbalancers": {
				Type:     schema.TypeList,
				Elem:     loadbalancersSchema(),
				Computed: true,
			},
		},
	}
}

func loadbalancersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cross_vpc_backend": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l7_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"autoscaling_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"min_l7_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backend_subnets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"protection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceElbLoadbalancersV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var (
		listLoadbalancersHttpUrl = "v3/{project_id}/elb/loadbalancers"
		listLoadbalancersProduct = "elb"
	)
	listLoadbalancersClient, err := cfg.NewServiceClient(listLoadbalancersProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listLoadbalancersPath := listLoadbalancersClient.Endpoint + listLoadbalancersHttpUrl
	listLoadbalancersPath = strings.ReplaceAll(listLoadbalancersPath, "{project_id}", listLoadbalancersClient.ProjectID)

	listLoadbalancersQueryParams := buildListLoadbalancersQueryParams(d)
	listLoadbalancersPath += listLoadbalancersQueryParams

	listLoadbalancersResp, err := pagination.ListAllItems(
		listLoadbalancersClient,
		"name",
		listLoadbalancersPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Loadbalancers")
	}

	listLoadbalancersRespJson, err := json.Marshal(listLoadbalancersResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listLoadbalancersRespBody interface{}
	err = json.Unmarshal(listLoadbalancersRespJson, &listLoadbalancersRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", cfg.GetRegion(d)),
		d.Set("loadbalancers", flattenListloadbalancersBody(listLoadbalancersRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListLoadbalancersQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("loadbalancer_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("share_type"); ok {
		if v == "dedicated" {
			res = fmt.Sprintf("%s&guaranteed=%v", res, "true")
		}
		if v == "share" {
			res = fmt.Sprintf("%s&guaranteed=%v", res, "false")
		}
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s&vpc_id=%v", res, v)
	}
	if v, ok := d.GetOk("ipv4_subnet_id"); ok {
		res = fmt.Sprintf("%s&vip_subnet_cidr_id=%v", res, v)
	}
	if v, ok := d.GetOk("ipv6_subnet_id"); ok {
		res = fmt.Sprintf("%s&ipv6_vip_virsubnet_id=%v", res, v)
	}
	if v, ok := d.GetOk("l4_flavor_id"); ok {
		res = fmt.Sprintf("%s&l4_flavor_id=%v", res, v)
	}
	if v, ok := d.GetOk("l7_flavor_id"); ok {
		res = fmt.Sprintf("%s&l7_flavor_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListloadbalancersBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("pools", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"availability_zone":     utils.PathSearch("availability_zone_list", v, nil),
			"cross_vpc_backend":     utils.PathSearch("ip_target_enable", v, nil),
			"vpc_id":                utils.PathSearch("vpc_id", v, nil),
			"ipv4_subnet_id":        utils.PathSearch("vip_subnet_cidr_id", v, nil),
			"ipv6_network_id":       utils.PathSearch("ipv6_vip_virsubnet_id", v, nil),
			"ipv4_address":          utils.PathSearch("vip_address", v, nil),
			"ipv4_port_id":          utils.PathSearch("vip_port_id", v, nil),
			"ipv6_address":          utils.PathSearch("ipv6_vip_address", v, nil),
			"l4_flavor_id":          utils.PathSearch("l4_flavor_id", v, nil),
			"l7_flavor_id":          utils.PathSearch("l7_flavor_id", v, nil),
			"region":                utils.PathSearch("region", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"autoscaling_enabled":   utils.PathSearch("enable", v, nil),
			"backend_subnets":       utils.PathSearch("elb_virsubnet_ids", v, nil),
			"protection_status":     utils.PathSearch("protection_status", v, nil),
			"protection_reason":     utils.PathSearch("protection_reason", v, nil),
		})
	}
	return rst
}
