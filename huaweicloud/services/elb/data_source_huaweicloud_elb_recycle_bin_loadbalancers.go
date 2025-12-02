package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/recycle-bin/loadbalancers
func DataSourceElbRecycleBinLoadBalancers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbRecycleBinLoadBalancersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"operating_status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"guaranteed": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vip_port_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vip_subnet_cidr_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_vip_port_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_vip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_vip_virsubnet_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"publicips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"availability_zone_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"l4_flavor_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"l7_flavor_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"billing_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"member_device_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"member_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_version": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"deletion_protection_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"elb_virsubnet_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"protection_status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"global_eips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loadbalancers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinLoadBalancersSchema(),
			},
		},
	}
}

func recycleBinLoadBalancersSchema() *schema.Resource {
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
			"availability_zone_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_scale_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l7_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l7_scale_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_vip_virsubnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_vip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_target_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinLoadBalancerPoolsSchema(),
			},
			"global_eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinLoadBalancerGlobalEipsSchema(),
			},
			"frozen_scene": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_bandwidth": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinLoadBalancerIpv6BandwidthSchema(),
			},
			"provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinLoadBalancerPublicipsSchema(),
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_protection_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"elb_virsubnet_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"waf_failure_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_vip_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guaranteed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"billing_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"elb_virsubnet_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinLoadBalancerListenersSchema(),
			},
			"vip_subnet_cidr_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinLoadBalancerTagsSchema(),
			},
			"auto_terminate_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func recycleBinLoadBalancerPoolsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func recycleBinLoadBalancerGlobalEipsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"global_eip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_eip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func recycleBinLoadBalancerIpv6BandwidthSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func recycleBinLoadBalancerPublicipsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"publicip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func recycleBinLoadBalancerListenersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func recycleBinLoadBalancerTagsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceElbRecycleBinLoadBalancersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/elb/recycle-bin/loadbalancers"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listQueryParams := buildListRecycleBinLoadBalancersQueryParams(d)
	listPath += listQueryParams

	listResp, err := pagination.ListAllItems(
		client,
		"marker",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving ELB recycle bin loadbalancers: %s", err)
	}
	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("loadbalancers", flattenRecycleBinLoadBalancersBody(listRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListRecycleBinLoadBalancersQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("loadbalancer_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "id"))
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "name"))
		for _, id := range v.([]interface{}) {
			res = fmt.Sprintf("%s&name=%v", res, id)
		}
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "description"))
	}
	if v, ok := d.GetOk("operating_status"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "operating_status"))
	}
	if v, ok := d.GetOk("guaranteed"); ok {
		guaranteed, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&guaranteed=%v", res, guaranteed)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "vpc_id"))
	}
	if v, ok := d.GetOk("vip_port_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "vip_port_id"))
	}
	if v, ok := d.GetOk("vip_address"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "vip_address"))
	}
	if v, ok := d.GetOk("vip_subnet_cidr_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "vip_subnet_cidr_id"))
	}
	if v, ok := d.GetOk("ipv6_vip_port_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "ipv6_vip_port_id"))
	}
	if v, ok := d.GetOk("ipv6_vip_address"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "ipv6_vip_address"))
	}
	if v, ok := d.GetOk("ipv6_vip_virsubnet_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "ipv6_vip_virsubnet_id"))
	}
	if v, ok := d.GetOk("publicips"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "publicips"))
	}
	if v, ok := d.GetOk("availability_zone_list"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "availability_zone_list"))
	}
	if v, ok := d.GetOk("l4_flavor_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "l4_flavor_id"))
	}
	if v, ok := d.GetOk("l7_flavor_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "l7_flavor_id"))
	}
	if v, ok := d.GetOk("billing_info"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "billing_info"))
	}
	if v, ok := d.GetOk("member_device_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "member_device_id"))
	}
	if v, ok := d.GetOk("member_address"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "member_address"))
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "enterprise_project_id"))
	}
	if v, ok := d.GetOk("ip_version"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "ip_version"))
	}
	if v, ok := d.GetOk("deletion_protection_enable"); ok {
		deletionProtectionEnable, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&deletion_protection_enable=%v", res, deletionProtectionEnable)
	}
	if v, ok := d.GetOk("elb_virsubnet_type"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "elb_virsubnet_type"))
	}
	if v, ok := d.GetOk("protection_status"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "protection_status"))
	}
	if v, ok := d.GetOk("global_eips"); ok {
		res = fmt.Sprintf("%s%v", res, buildCycleParam(v.([]interface{}), "global_eips"))
	}
	if v, ok := d.GetOk("log_topic_id"); ok {
		res = fmt.Sprintf("%s&log_topic_id=%v", res, v)
	}
	if v, ok := d.GetOk("log_group_id"); ok {
		res = fmt.Sprintf("%s&log_group_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func buildCycleParam(value []interface{}, param string) string {
	res := ""
	for _, v := range value {
		res = fmt.Sprintf("%s&%s=%v", res, param, v)
	}
	return res
}

func flattenRecycleBinLoadBalancersBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("loadbalancers", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                         utils.PathSearch("id", v, nil),
			"name":                       utils.PathSearch("name", v, nil),
			"availability_zone_list":     utils.PathSearch("availability_zone_list", v, nil),
			"description":                utils.PathSearch("description", v, nil),
			"vpc_id":                     utils.PathSearch("vpc_id", v, nil),
			"l4_flavor_id":               utils.PathSearch("l4_flavor_id", v, nil),
			"l4_scale_flavor_id":         utils.PathSearch("l4_scale_flavor_id", v, nil),
			"l7_flavor_id":               utils.PathSearch("l7_flavor_id", v, nil),
			"l7_scale_flavor_id":         utils.PathSearch("l4_scale_flavor_id", v, nil),
			"ipv6_vip_virsubnet_id":      utils.PathSearch("ipv6_vip_virsubnet_id", v, nil),
			"ipv6_vip_address":           utils.PathSearch("ipv6_vip_address", v, nil),
			"ip_target_enable":           utils.PathSearch("ip_target_enable", v, nil),
			"pools":                      flattenRecycleBinLoadBalancerPoolsBody(v),
			"global_eips":                flattenRecycleBinLoadBalancerGlobalEipsBody(v),
			"frozen_scene":               utils.PathSearch("frozen_scene", v, nil),
			"ipv6_bandwidth":             flattenRecycleBinLoadBalancerIpv6BandwidthBody(v),
			"provider":                   utils.PathSearch("provider", v, nil),
			"protection_status":          utils.PathSearch("protection_status", v, nil),
			"vip_address":                utils.PathSearch("vip_address", v, nil),
			"vip_port_id":                utils.PathSearch("vip_port_id", v, nil),
			"publicips":                  flattenRecycleBinLoadBalancerPublicipsBody(v),
			"charge_mode":                utils.PathSearch("charge_mode", v, nil),
			"operating_status":           utils.PathSearch("operating_status", v, nil),
			"enterprise_project_id":      utils.PathSearch("enterprise_project_id", v, nil),
			"deletion_protection_enable": utils.PathSearch("deletion_protection_enable", v, nil),
			"provisioning_status":        utils.PathSearch("provisioning_status", v, nil),
			"elb_virsubnet_ids":          utils.PathSearch("elb_virsubnet_ids", v, nil),
			"public_border_group":        utils.PathSearch("public_border_group", v, nil),
			"waf_failure_action":         utils.PathSearch("waf_failure_action", v, nil),
			"ipv6_vip_port_id":           utils.PathSearch("ipv6_vip_port_id", v, nil),
			"guaranteed":                 utils.PathSearch("guaranteed", v, nil),
			"billing_info":               utils.PathSearch("billing_info", v, nil),
			"elb_virsubnet_type":         utils.PathSearch("elb_virsubnet_type", v, nil),
			"protection_reason":          utils.PathSearch("protection_reason", v, nil),
			"log_topic_id":               utils.PathSearch("log_topic_id", v, nil),
			"listeners":                  flattenRecycleBinLoadBalancerListenersBody(v),
			"vip_subnet_cidr_id":         utils.PathSearch("vip_subnet_cidr_id", v, nil),
			"tags":                       flattenRecycleBinLoadBalancerTagsBody(v),
			"auto_terminate_time":        utils.PathSearch("auto_terminate_time", v, nil),
			"created_at":                 utils.PathSearch("created_at", v, nil),
			"updated_at":                 utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenRecycleBinLoadBalancerPoolsBody(resp interface{}) []interface{} {
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
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenRecycleBinLoadBalancerGlobalEipsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("global_eips", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"global_eip_id":      utils.PathSearch("global_eip_id", v, nil),
			"global_eip_address": utils.PathSearch("global_eip_address", v, nil),
			"ip_version":         utils.PathSearch("ip_version", v, nil),
		})
	}
	return rst
}

func flattenRecycleBinLoadBalancerIpv6BandwidthBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("ipv6_bandwidth", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenRecycleBinLoadBalancerPublicipsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("publicips", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"publicip_id":      utils.PathSearch("publicip_id", v, nil),
			"publicip_address": utils.PathSearch("publicip_address", v, nil),
			"ip_version":       utils.PathSearch("ip_version", v, nil),
		})
	}
	return rst
}

func flattenRecycleBinLoadBalancerListenersBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("listeners", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenRecycleBinLoadBalancerTagsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("tags", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
