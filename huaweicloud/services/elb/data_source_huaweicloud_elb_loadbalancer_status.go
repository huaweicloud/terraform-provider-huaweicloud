package elb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/statuses
func DataSourceElbLoadBalancerStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbLoadBalancerStatusRead,

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
			"statuses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesSchema(),
			},
		},
	}
}

func loadBalancerStatusesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"loadbalancer": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesLoadBalancerSchema(),
			},
		},
	}
	return &sc
}

func loadBalancerStatusesLoadBalancerSchema() *schema.Resource {
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
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesLoadBalancerListenerSchema(),
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesLoadBalancerPoolSchema(),
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func loadBalancerStatusesLoadBalancerListenerSchema() *schema.Resource {
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
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesLoadBalancerPoolSchema(),
			},
			"l7policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesLoadBalancerL7PolicySchema(),
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func loadBalancerStatusesLoadBalancerPoolSchema() *schema.Resource {
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
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthmonitor": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesLoadBalancerPoolHealthMonitorSchema(),
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesLoadBalancerMemberSchema(),
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func loadBalancerStatusesLoadBalancerPoolHealthMonitorSchema() *schema.Resource {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func loadBalancerStatusesLoadBalancerMemberSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func loadBalancerStatusesLoadBalancerL7PolicySchema() *schema.Resource {
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
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerStatusesLoadBalancerL7PolicyRuleSchema(),
			},
		},
	}
	return &sc
}

func loadBalancerStatusesLoadBalancerL7PolicyRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceElbLoadBalancerStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/statuses"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{loadbalancer_id}", d.Get("loadbalancer_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving ELB load balancer status: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("statuses", flattenLoadBalancerStatusBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLoadBalancerStatusBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("statuses", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"loadbalancer": flattenLoadBalancerStatusLoadBalancerBody(curJson),
		},
	}
	return rst
}

func flattenLoadBalancerStatusLoadBalancerBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("loadbalancer", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"id":                  utils.PathSearch("id", curJson, nil),
			"name":                utils.PathSearch("name", curJson, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", curJson, nil),
			"listeners":           flattenLoadBalancerStatusLoadBalancerListenerBody(curJson),
			"pools":               flattenLoadBalancerStatusLoadBalancerPoolBody(curJson),
			"operating_status":    utils.PathSearch("operating_status", curJson, nil),
		},
	}
	return rst
}

func flattenLoadBalancerStatusLoadBalancerListenerBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("listeners", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"name":                utils.PathSearch("name", v, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", v, nil),
			"pools":               flattenLoadBalancerStatusLoadBalancerPoolBody(v),
			"l7policies":          flattenLoadBalancerStatusLoadBalancerListenerL7PolicierBody(v),
			"operating_status":    utils.PathSearch("operating_status", v, nil),
		})
	}
	return rst
}

func flattenLoadBalancerStatusLoadBalancerPoolBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("pools", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"name":                utils.PathSearch("name", v, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", v, nil),
			"healthmonitor":       flattenLoadBalancerStatusLoadBalancerPoolHealthMonitorBody(v),
			"members":             flattenLoadBalancerStatusLoadBalancerPoolMemberBody(v),
			"operating_status":    utils.PathSearch("operating_status", v, nil),
		})
	}
	return rst
}

func flattenLoadBalancerStatusLoadBalancerPoolHealthMonitorBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("healthmonitor", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"id":                  utils.PathSearch("id", curJson, nil),
			"name":                utils.PathSearch("name", curJson, nil),
			"type":                utils.PathSearch("type", curJson, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", curJson, nil),
		},
	}
	return rst
}

func flattenLoadBalancerStatusLoadBalancerPoolMemberBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("members", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"address":             utils.PathSearch("name", v, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", v, nil),
			"protocol_port":       utils.PathSearch("protocol_port", v, nil),
			"operating_status":    utils.PathSearch("operating_status", v, nil),
		})
	}
	return rst
}

func flattenLoadBalancerStatusLoadBalancerListenerL7PolicierBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("l7policies", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"name":                utils.PathSearch("name", v, nil),
			"action":              utils.PathSearch("action", v, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", v, nil),
			"rules":               flattenLoadBalancerStatusLoadBalancerListenerL7PolicyRuleBody(v),
		})
	}
	return rst
}

func flattenLoadBalancerStatusLoadBalancerListenerL7PolicyRuleBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("rules", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"type":                utils.PathSearch("type", v, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", v, nil),
		})
	}
	return rst
}
