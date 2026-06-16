package elb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/topology
func DataSourceElbLoadBalancerTopology() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbLoadBalancerTopologyRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pool_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyNodesSchema(),
			},
			"edges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyEdgesSchema(),
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyLabelsSchema(),
			},
		},
	}
}

func topologyNodesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"loadbalancers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyLoadBalancerNodeSchema(),
			},
			"eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyEipNodeSchema(),
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyListenerNodeSchema(),
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyPoolNodeSchema(),
			},
		},
	}
}

func topologyLoadBalancerNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guaranteed": {
				Type:     schema.TypeBool,
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
			"vip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_vip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func topologyEipNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func topologyListenerNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyPortRangeSchema(),
			},
		},
	}
}

func topologyPortRangeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"end_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func topologyPoolNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func topologyEdgesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func topologyLabelsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"l7policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyPolicyLabelSchema(),
			},
		},
	}
}

func topologyPolicyLabelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// The type in API doucument is `string`, actually return `int`
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyL7RulesSchema(),
			},
		},
	}
}

func topologyL7RulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"invert": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyL7RuleConditionSchema(),
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
}

func topologyL7RuleConditionSchema() *schema.Resource {
	return &schema.Resource{
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
}

func dataSourceElbLoadBalancerTopologyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	httpUrl := "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/topology"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{loadbalancer_id}", d.Get("loadbalancer_id").(string))
	getPath += buildGetElbLoadBalancerTopologyQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving ELB load balancer topology: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("nodes", flattenElbLoadBalancerTopologyNodes(getRespBody)),
		d.Set("edges", flattenElbLoadBalancerTopologyEdges(getRespBody)),
		d.Set("labels", flattenElbLoadBalancerTopologyLabels(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetElbLoadBalancerTopologyQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&listener_id=%v", res, v)
	}
	if v, ok := d.GetOk("pool_id"); ok {
		res = fmt.Sprintf("%s&pool_id=%v", res, v)
	}
	if v, ok := d.GetOk("listener_name"); ok {
		res = fmt.Sprintf("%s&listener_name=%v", res, v)
	}
	if v, ok := d.GetOk("listener_protocol"); ok {
		res = fmt.Sprintf("%s&listener_protocol=%v", res, v)
	}
	if v, ok := d.GetOk("listener_protocol_port"); ok {
		res = fmt.Sprintf("%s&listener_protocol_port=%v", res, v)
	}
	if v, ok := d.GetOk("pool_name"); ok {
		res = fmt.Sprintf("%s&pool_name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenElbLoadBalancerTopologyNodes(resp interface{}) []interface{} {
	curJson := utils.PathSearch("topology.nodes", resp, nil)
	if curJson == nil {
		return nil
	}

	res := []interface{}{
		map[string]interface{}{
			"loadbalancers": flattenElbLoadBalancerTopologyLoadBalancers(curJson),
			"eips":          flattenElbLoadBalancerTopologyEips(curJson),
			"listeners":     flattenElbLoadBalancerTopologyListeners(curJson),
			"pools":         flattenElbLoadBalancerTopologyPools(curJson),
		},
	}
	return res
}

func flattenElbLoadBalancerTopologyLoadBalancers(resp interface{}) []interface{} {
	curJson := utils.PathSearch("loadbalancers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"guaranteed":             utils.PathSearch("guaranteed", v, nil),
			"l4_flavor_id":           utils.PathSearch("l4_flavor_id", v, nil),
			"l7_flavor_id":           utils.PathSearch("l7_flavor_id", v, nil),
			"vip_address":            utils.PathSearch("vip_address", v, nil),
			"ipv6_vip_address":       utils.PathSearch("ipv6_vip_address", v, nil),
			"availability_zone_list": utils.PathSearch("availability_zone_list", v, nil),
		})
	}
	return res
}

func flattenElbLoadBalancerTopologyEips(resp interface{}) []interface{} {
	curJson := utils.PathSearch("eips", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"ip_address": utils.PathSearch("ip_address", v, nil),
			"ip_version": utils.PathSearch("ip_version", v, nil),
		})
	}
	return res
}

func flattenElbLoadBalancerTopologyListeners(resp interface{}) []interface{} {
	curJson := utils.PathSearch("listeners", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"protocol":      utils.PathSearch("protocol", v, nil),
			"protocol_port": utils.PathSearch("protocol_port", v, nil),
			"port_ranges":   flattenElbLoadBalancerTopologyPortRanges(v),
		})
	}
	return res
}

func flattenElbLoadBalancerTopologyPortRanges(resp interface{}) []interface{} {
	curJson := utils.PathSearch("port_ranges", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"start_port": utils.PathSearch("start_port", v, nil),
			"end_port":   utils.PathSearch("end_port", v, nil),
		})
	}
	return res
}

func flattenElbLoadBalancerTopologyPools(resp interface{}) []interface{} {
	curJson := utils.PathSearch("pools", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"name":         utils.PathSearch("name", v, nil),
			"protocol":     utils.PathSearch("protocol", v, nil),
			"type":         utils.PathSearch("type", v, nil),
			"lb_algorithm": utils.PathSearch("lb_algorithm", v, nil),
		})
	}
	return res
}

func flattenElbLoadBalancerTopologyEdges(resp interface{}) []interface{} {
	curJson := utils.PathSearch("topology.edges", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"source":      utils.PathSearch("source", v, nil),
			"target":      utils.PathSearch("target", v, nil),
			"source_type": utils.PathSearch("source_type", v, nil),
			"target_type": utils.PathSearch("target_type", v, nil),
			"label":       utils.PathSearch("label", v, nil),
			"label_id":    utils.PathSearch("label_id", v, nil),
		})
	}
	return res
}

func flattenElbLoadBalancerTopologyLabels(resp interface{}) []interface{} {
	curJson := utils.PathSearch("topology.labels", resp, nil)
	if curJson == nil {
		return nil
	}

	res := []interface{}{
		map[string]interface{}{
			"l7policies": flattenElbLoadBalancerTopologyL7Policies(curJson),
		},
	}
	return res
}

func flattenElbLoadBalancerTopologyL7Policies(resp interface{}) []interface{} {
	curJson := utils.PathSearch("l7policies", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":       utils.PathSearch("id", v, nil),
			"name":     utils.PathSearch("name", v, nil),
			"priority": utils.PathSearch("priority", v, nil),
			"action":   utils.PathSearch("action", v, nil),
			"rules":    flattenTopologyL7Rules(v),
		})
	}
	return res
}

func flattenTopologyL7Rules(resp interface{}) []interface{} {
	curJson := utils.PathSearch("rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"type":                utils.PathSearch("type", v, nil),
			"compare_type":        utils.PathSearch("compare_type", v, nil),
			"key":                 utils.PathSearch("key", v, nil),
			"value":               utils.PathSearch("value", v, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", v, nil),
			"invert":              utils.PathSearch("invert", v, nil),
			"conditions":          flattenTopologyL7RulesConditions(v),
			"created_at":          utils.PathSearch("created_at", v, nil),
			"updated_at":          utils.PathSearch("updated_at", v, nil),
		})
	}
	return res
}

func flattenTopologyL7RulesConditions(resp interface{}) []interface{} {
	curJson := utils.PathSearch("conditions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return res
}
