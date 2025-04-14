package lb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LB GET /v2/{project_id}/elb/pools
func DataSourcePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourcePoolsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pool_id": {
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"healthmonitor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lb_method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pools": {
				Type:     schema.TypeList,
				Elem:     poolsPoolsSchema(),
				Computed: true,
			},
		},
	}
}

func poolsPoolsSchema() *schema.Resource {
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
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthmonitor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Elem:     poolsPoolListenersSchema(),
				Computed: true,
			},
			"loadbalancers": {
				Type:     schema.TypeList,
				Elem:     poolsPoolLoadbalancersSchema(),
				Computed: true,
			},
			"members": {
				Type:     schema.TypeList,
				Elem:     poolsPoolMembersSchema(),
				Computed: true,
			},
			"persistence": {
				Type:     schema.TypeList,
				Elem:     poolsPoolPersistenceSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func poolsPoolListenersSchema() *schema.Resource {
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

func poolsPoolLoadbalancersSchema() *schema.Resource {
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

func poolsPoolMembersSchema() *schema.Resource {
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

func poolsPoolPersistenceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cookie_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourcePoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listPools: Query the List of LB pools
	var (
		listPoolsHttpUrl = "v2/{project_id}/elb/pools"
		listPoolsProduct = "elb"
	)
	listPoolsClient, err := cfg.NewServiceClient(listPoolsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Pools Client: %s", err)
	}

	listPoolsPath := listPoolsClient.Endpoint + listPoolsHttpUrl
	listPoolsPath = strings.ReplaceAll(listPoolsPath, "{project_id}", listPoolsClient.ProjectID)

	listPoolsPath += buildListPoolsQueryParams(d)

	listPoolsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listPoolsResp, err := listPoolsClient.Request("GET", listPoolsPath, &listPoolsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving pools")
	}

	listPoolsRespBody, err := utils.FlattenResponse(listPoolsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("pools", flattenListPoolsBodyPools(listPoolsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListPoolsBodyPools(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("pools", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"protocol":          utils.PathSearch("protocol", v, nil),
			"lb_method":         utils.PathSearch("lb_algorithm", v, nil),
			"healthmonitor_id":  utils.PathSearch("healthmonitor_id", v, nil),
			"protection_status": utils.PathSearch("protection_status", v, nil),
			"protection_reason": utils.PathSearch("protection_reason", v, nil),
			"listeners":         flattenPoolListeners(v),
			"loadbalancers":     flattenPoolLoadbalancers(v),
			"members":           flattenPoolMembers(v),
			"persistence":       flattenPoolPersistence(v),
		})
	}
	return rst
}

func flattenPoolListeners(resp interface{}) []interface{} {
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

func flattenPoolLoadbalancers(resp interface{}) []interface{} {
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
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenPoolMembers(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("members", resp, make([]interface{}, 0))
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

func flattenPoolPersistence(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("session_persistence", resp, nil)
	if curJson == nil {
		return rst
	}
	if curJson == nil {
		return nil
	}

	rst = []interface{}{
		map[string]interface{}{
			"type":        utils.PathSearch("type", curJson, nil),
			"cookie_name": utils.PathSearch("cookie_name", curJson, nil),
			"timeout":     utils.PathSearch("persistence_timeout", curJson, nil),
		},
	}
	return rst
}

func buildListPoolsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("pool_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if v, ok := d.GetOk("loadbalancer_id"); ok {
		res = fmt.Sprintf("%s&loadbalancer_id=%v", res, v)
	}
	if v, ok := d.GetOk("member_address"); ok {
		res = fmt.Sprintf("%s&member_address=%v", res, v)
	}
	if v, ok := d.GetOk("member_device_id"); ok {
		res = fmt.Sprintf("%s&member_device_id=%v", res, v)
	}
	if v, ok := d.GetOk("healthmonitor_id"); ok {
		res = fmt.Sprintf("%s&healthmonitor_id=%v", res, v)
	}
	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&protocol=%v", res, v)
	}
	if v, ok := d.GetOk("lb_method"); ok {
		res = fmt.Sprintf("%s&lb_algorithm=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
