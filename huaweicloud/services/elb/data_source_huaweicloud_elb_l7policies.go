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

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/l7policies
func DataSourceElbL7policies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbL7policiesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"l7policy_id": {
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
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redirect_listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redirect_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"l7policies": {
				Type:     schema.TypeList,
				Elem:     l7policiesSchema(),
				Computed: true,
			},
		},
	}
}

func l7policiesSchema() *schema.Resource {
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"redirect_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"redirect_pools_config": {
				Type:     schema.TypeList,
				Elem:     redirectPoolsSchema(),
				Computed: true,
			},
			"redirect_listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
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
			"redirect_url_config": {
				Type:     schema.TypeList,
				Elem:     redirectUrlSchema(),
				Computed: true,
			},
			"fixed_response_config": {
				Type:     schema.TypeList,
				Elem:     fixedResponseSchema(),
				Computed: true,
			},
			"redirect_pools_extend_config": {
				Type:     schema.TypeList,
				Elem:     redirectPoolsExtendSchema(),
				Computed: true,
			},
			"redirect_pools_sticky_session_config": {
				Type:     schema.TypeList,
				Elem:     redirectPoolsStickySessionSchema(),
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

func redirectPoolsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func redirectUrlSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insert_headers_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     insertHeadersSchema(),
			},
			"remove_headers_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     removeHeadersSchema(),
			},
		},
	}
	return &sc
}

func fixedResponseSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insert_headers_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     insertHeadersSchema(),
			},
			"remove_headers_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     removeHeadersSchema(),
			},
			"traffic_limit_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     trafficLimitSchema(),
			},
		},
	}
	return &sc
}

func redirectPoolsExtendSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rewrite_url_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rewrite_url_config": {
				Type:     schema.TypeList,
				Elem:     rewriteUrlSchema(),
				Computed: true,
			},
			"insert_headers_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     insertHeadersSchema(),
			},
			"remove_headers_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     removeHeadersSchema(),
			},
			"traffic_limit_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     trafficLimitSchema(),
			},
		},
	}
	return &sc
}

func rewriteUrlSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func insertHeadersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     insertHeaderSchema(),
			},
		},
	}
	return &sc
}

func insertHeaderSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_type": {
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

func removeHeadersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     removeHeaderSchema(),
			},
		},
	}
	return &sc
}

func removeHeaderSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func trafficLimitSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"qps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"per_source_ip_qps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"burst": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func redirectPoolsStickySessionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:     schema.TypeBool,
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

func dataSourceElbL7policiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		listL7policiesHttpUrl = "v3/{project_id}/elb/l7policies"
		listL7policiesProduct = "elb"
	)
	listL7policiesClient, err := cfg.NewServiceClient(listL7policiesProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listL7policiesPath := listL7policiesClient.Endpoint + listL7policiesHttpUrl
	listL7policiesPath = strings.ReplaceAll(listL7policiesPath, "{project_id}", listL7policiesClient.ProjectID)
	listL7policiesQueryParams := buildListL7policiesQueryParams(d)
	listL7policiesPath += listL7policiesQueryParams
	listL7policiesResp, err := pagination.ListAllItems(
		listL7policiesClient,
		"marker",
		listL7policiesPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB L7 policies")
	}

	listL7policiesRespJson, err := json.Marshal(listL7policiesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listL7policiesRespBody interface{}
	err = json.Unmarshal(listL7policiesRespJson, &listL7policiesRespBody)
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
		d.Set("l7policies", flattenListL7policiesBody(listL7policiesRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListL7policiesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("l7policy_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("display_all_rules"); ok {
		displayAllRules, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&display_all_rules=%v", res, displayAllRules)
	}
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&listener_id=%v", res, v)
	}
	if v, ok := d.GetOk("action"); ok {
		res = fmt.Sprintf("%s&action=%v", res, v)
	}
	if v, ok := d.GetOk("priority"); ok {
		res = fmt.Sprintf("%s&priority=%v", res, v)
	}
	if v, ok := d.GetOk("provisioning_status"); ok {
		res = fmt.Sprintf("%s&provisioning_status=%v", res, v)
	}
	if v, ok := d.GetOk("redirect_listener_id"); ok {
		res = fmt.Sprintf("%s&redirect_listener_id=%v", res, v)
	}
	if v, ok := d.GetOk("redirect_pool_id"); ok {
		res = fmt.Sprintf("%s&redirect_pool_id=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListL7policiesBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("l7policies", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                                   utils.PathSearch("id", v, nil),
			"name":                                 utils.PathSearch("name", v, nil),
			"description":                          utils.PathSearch("description", v, nil),
			"enterprise_project_id":                utils.PathSearch("enterprise_project_id", v, nil),
			"priority":                             utils.PathSearch("priority", v, nil),
			"provisioning_status":                  utils.PathSearch("provisioning_status", v, nil),
			"listener_id":                          utils.PathSearch("listener_id", v, nil),
			"redirect_pool_id":                     utils.PathSearch("redirect_pool_id", v, nil),
			"redirect_pools_config":                flattenPoolsConfigBody(v),
			"redirect_listener_id":                 utils.PathSearch("redirect_listener_id", v, nil),
			"action":                               utils.PathSearch("action", v, nil),
			"rules":                                utils.PathSearch("rules", v, nil),
			"redirect_url_config":                  flattenRedirectUrlConfigBody(v),
			"fixed_response_config":                flattenFixedResponseConfigBody(v),
			"redirect_pools_extend_config":         flattenRedirectPoolsExtendConfigBody(v),
			"redirect_pools_sticky_session_config": flattenRedirectPoolsStickySessionConfigBody(v),
			"created_at":                           utils.PathSearch("created_at", v, nil),
			"updated_at":                           utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenPoolsConfigBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("redirect_pools_config", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"pool_id": utils.PathSearch("pool_id", curJson, nil),
			"weight":  utils.PathSearch("weight", curJson, nil),
		},
	}
	return rst
}

func flattenRedirectUrlConfigBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("redirect_url_config", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"status_code":           utils.PathSearch("status_code", curJson, nil),
			"protocol":              utils.PathSearch("protocol", curJson, nil),
			"host":                  utils.PathSearch("host", curJson, nil),
			"port":                  utils.PathSearch("port", curJson, nil),
			"path":                  utils.PathSearch("path", curJson, nil),
			"query":                 utils.PathSearch("query", curJson, nil),
			"insert_headers_config": flattenInsertHeadersConfigBody(curJson),
			"remove_headers_config": flattenRemoveHeadersConfigBody(curJson),
		},
	}
	return rst
}

func flattenInsertHeadersConfigBody(cfg interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("insert_headers_config", cfg, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"configs": flattenInsertHeaderConfigsBody(curJson),
		},
	}
	return rst
}

func flattenInsertHeaderConfigsBody(insertHeaderConfigs interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs", insertHeaderConfigs, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":        utils.PathSearch("key", v, nil),
			"value":      utils.PathSearch("value", v, nil),
			"value_type": utils.PathSearch("value_type", v, nil),
		})
	}
	return rst
}

func flattenRemoveHeadersConfigBody(cfg interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("remove_headers_config", cfg, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"configs": flattenRemoveHeaderConfigsBody(curJson),
		},
	}
	return rst
}

func flattenRemoveHeaderConfigsBody(removeHeaderConfigs interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs", removeHeaderConfigs, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key": utils.PathSearch("key", v, nil),
		})
	}
	return rst
}

func flattenFixedResponseConfigBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("fixed_response_config", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"status_code":           utils.PathSearch("status_code", curJson, nil),
			"content_type":          utils.PathSearch("content_type", curJson, nil),
			"message_body":          utils.PathSearch("message_body", curJson, nil),
			"insert_headers_config": flattenInsertHeadersConfigBody(curJson),
			"remove_headers_config": flattenRemoveHeadersConfigBody(curJson),
			"traffic_limit_config":  flattenTrafficLimitConfigBody(curJson),
		},
	}
	return rst
}

func flattenTrafficLimitConfigBody(cfg interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("traffic_limit_config", cfg, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"qps":               utils.PathSearch("qps", curJson, nil),
			"per_source_ip_qps": utils.PathSearch("per_source_ip_qps", curJson, nil),
			"burst":             utils.PathSearch("burst", curJson, nil),
		},
	}
	return rst
}

func flattenRedirectPoolsExtendConfigBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("redirect_pools_extend_config", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"rewrite_url_enabled":   utils.PathSearch("rewrite_url_enabled", curJson, nil),
			"rewrite_url_config":    flattenRewriteUrlConfigBody(curJson),
			"insert_headers_config": flattenInsertHeadersConfigBody(curJson),
			"remove_headers_config": flattenRemoveHeadersConfigBody(curJson),
			"traffic_limit_config":  flattenTrafficLimitConfigBody(curJson),
		},
	}
	return rst
}

func flattenRewriteUrlConfigBody(cfg interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("rewrite_url_config", cfg, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"host":  utils.PathSearch("host", curJson, nil),
			"path":  utils.PathSearch("path", curJson, nil),
			"query": utils.PathSearch("query", curJson, nil),
		},
	}
	return rst
}

func flattenRedirectPoolsStickySessionConfigBody(resp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("redirect_pools_sticky_session_config", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"enable":  utils.PathSearch("enable", curJson, nil),
			"timeout": utils.PathSearch("timeout", curJson, nil),
		},
	}
	return rst
}
