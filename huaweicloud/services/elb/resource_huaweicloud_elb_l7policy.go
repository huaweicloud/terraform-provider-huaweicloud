package elb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB POST /v3/{project_id}/elb/l7policies
// @API ELB GET /v3/{project_id}/elb/l7policies/{l7policy_id}
// @API ELB PUT /v3/{project_id}/elb/l7policies/{l7policy_id}
// @API ELB DELETE /v3/{project_id}/elb/l7policies/{l7policy_id}
func ResourceL7PolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7PolicyV3Create,
		ReadContext:   resourceL7PolicyV3Read,
		UpdateContext: resourceL7PolicyV3Update,
		DeleteContext: resourceL7PolicyV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "REDIRECT_TO_POOL",
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"redirect_listener_id": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{"redirect_listener_id", "redirect_pool_id", "redirect_pools_config", "redirect_url_config",
					"fixed_response_config"},
				Computed: true,
			},
			"redirect_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"redirect_pools_config": {
				Type:          schema.TypeSet,
				Elem:          redirectPoolsConfigSchema(),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"redirect_pool_id"},
			},
			"redirect_pools_sticky_session_config": {
				Type:          schema.TypeList,
				Elem:          redirectPoolsStickySessionConfigSchema(),
				MaxItems:      1,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"redirect_pool_id"},
			},
			"redirect_pools_extend_config": {
				Type:     schema.TypeList,
				Elem:     redirectPoolsExtendConfigSchema(),
				MaxItems: 1,
				Optional: true,
				Computed: true,
			},
			"redirect_url_config": {
				Type:     schema.TypeList,
				Elem:     redirectUrlConfigSchema(),
				MaxItems: 1,
				Optional: true,
				Computed: true,
			},
			"fixed_response_config": {
				Type:     schema.TypeList,
				Elem:     fixedResponseConfigSchema(),
				MaxItems: 1,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
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
}

func redirectPoolsExtendConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rewrite_url_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"rewrite_url_config": {
				Type:     schema.TypeList,
				Elem:     rewriteUrlConfigSchema(),
				MaxItems: 1,
				Optional: true,
				Computed: true,
			},
			"insert_headers_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     insertHeadersConfigSchema(),
				MaxItems: 1,
			},
			"remove_headers_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     removeHeadersConfigSchema(),
				MaxItems: 1,
			},
			"traffic_limit_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     trafficLimitConfigSchema(),
				MaxItems: 1,
			},
		},
	}
	return &sc
}

func redirectPoolsConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func redirectPoolsStickySessionConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func rewriteUrlConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	return &sc
}

func redirectUrlConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insert_headers_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     insertHeadersConfigSchema(),
				MaxItems: 1,
			},
			"remove_headers_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     removeHeadersConfigSchema(),
				MaxItems: 1,
			},
		},
	}
	return &sc
}

func fixedResponseConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"message_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insert_headers_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     insertHeadersConfigSchema(),
				MaxItems: 1,
			},
			"remove_headers_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     removeHeadersConfigSchema(),
				MaxItems: 1,
			},
			"traffic_limit_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     trafficLimitConfigSchema(),
				MaxItems: 1,
			},
		},
	}
	return &sc
}

func insertHeadersConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"configs": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     insertHeaderConfigSchema(),
			},
		},
	}
	return &sc
}

func insertHeaderConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func removeHeadersConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"configs": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     removeHeaderConfigSchema(),
			},
		},
	}
	return &sc
}

func removeHeaderConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func trafficLimitConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"qps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"per_source_ip_qps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"burst": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
	return &sc
}

func resourceL7PolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/l7policies"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateL7PolicyBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB L7 policy: %s", err)
	}

	createARespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error retrieving ELB L7 policy: %s", err)
	}
	policyId := utils.PathSearch("l7policy.id", createARespBody, "").(string)
	if policyId == "" {
		return diag.Errorf("error creating ELB L7 policy: ID is not found in API response")
	}

	d.SetId(policyId)

	// Wait for L7 Policy to become active before continuing
	err = waitForL7Policy(ctx, client, policyId, "ACTIVE", nil, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7PolicyV3Read(ctx, d, meta)
}
func buildCreateL7PolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                                 utils.ValueIgnoreEmpty(d.Get("name")),
		"description":                          utils.ValueIgnoreEmpty(d.Get("description")),
		"action":                               d.Get("action"),
		"priority":                             utils.ValueIgnoreEmpty(d.Get("priority")),
		"listener_id":                          d.Get("listener_id"),
		"redirect_pool_id":                     utils.ValueIgnoreEmpty(d.Get("redirect_pool_id")),
		"redirect_pools_config":                buildRedirectPoolsConfig(d),
		"redirect_pools_sticky_session_config": buildRedirectPoolsStickySessionConfig(d),
		"redirect_pools_extend_config":         buildRedirectPoolsExtendConfig(d),
		"redirect_listener_id":                 utils.ValueIgnoreEmpty(d.Get("redirect_listener_id")),
		"redirect_url_config":                  buildRedirectUrlConfig(d),
		"fixed_response_config":                buildFixedResponseConfig(d),
	}
	return map[string]interface{}{"l7policy": bodyParams}
}

func buildRedirectPoolsConfig(d *schema.ResourceData) []map[string]interface{} {
	rawRedirectPoolsConfig := d.Get("redirect_pools_config").(*schema.Set)
	if rawRedirectPoolsConfig.Len() == 0 {
		return nil
	}

	redirectPoolsConfig := make([]map[string]interface{}, 0, rawRedirectPoolsConfig.Len())
	for _, rawConfig := range rawRedirectPoolsConfig.List() {
		if v, ok := rawConfig.(map[string]interface{}); ok {
			redirectPoolsConfig = append(redirectPoolsConfig, map[string]interface{}{
				"pool_id": v["pool_id"],
				"weight":  utils.ValueIgnoreEmpty(v["weight"]),
			})
		}
	}
	return redirectPoolsConfig
}

func buildRedirectPoolsStickySessionConfig(d *schema.ResourceData) map[string]interface{} {
	if rawConfig, ok := d.GetOk("redirect_pools_sticky_session_config"); ok {
		if v, ok := rawConfig.([]interface{})[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"enable":  v["enable"].(bool),
				"timeout": v["timeout"].(int),
			}
			return params
		}
	}
	return nil
}

func buildRedirectPoolsExtendConfig(d *schema.ResourceData) map[string]interface{} {
	if rawConfig, ok := d.GetOk("redirect_pools_extend_config"); ok {
		if v, ok := rawConfig.([]interface{})[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"rewrite_url_enable":    v["rewrite_url_enabled"].(bool),
				"rewrite_url_config":    buildRewriteUrlConfig(v["rewrite_url_config"]),
				"insert_headers_config": buildInsertHeadersConfig(v["insert_headers_config"]),
				"remove_headers_config": buildRemoveHeadersConfig(v["remove_headers_config"]),
				"traffic_limit_config":  buildTrafficLimitConfig(v["traffic_limit_config"]),
			}
			return params
		}
	}
	return nil
}

func buildRewriteUrlConfig(data interface{}) map[string]interface{} {
	rewriteUrlConfigRaw := data.([]interface{})
	if len(rewriteUrlConfigRaw) == 1 {
		if v, ok := rewriteUrlConfigRaw[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"host":  v["host"].(string),
				"path":  v["path"].(string),
				"query": v["query"].(string),
			}
			return params
		}
	}
	return nil
}

func buildInsertHeadersConfig(data interface{}) map[string]interface{} {
	insertHeadersConfigRaw := data.([]interface{})
	if len(insertHeadersConfigRaw) == 1 {
		if v, ok := insertHeadersConfigRaw[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"configs": buildInsertHeaderConfig(v["configs"]),
			}
			return params
		}
	}
	return nil
}

func buildInsertHeaderConfig(data interface{}) []map[string]interface{} {
	insertHeaderConfigsRaw := data.(*schema.Set)
	if insertHeaderConfigsRaw.Len() == 0 {
		return nil
	}

	insertHeaderConfigs := make([]map[string]interface{}, 0, insertHeaderConfigsRaw.Len())
	for _, rawConfig := range insertHeaderConfigsRaw.List() {
		if v, ok := rawConfig.(map[string]interface{}); ok {
			insertHeaderConfigs = append(insertHeaderConfigs, map[string]interface{}{
				"key":        v["key"].(string),
				"value":      v["value"].(string),
				"value_type": v["value_type"].(string),
			})
		}
	}

	return insertHeaderConfigs
}

func buildRemoveHeadersConfig(data interface{}) map[string]interface{} {
	removeHeadersConfigRaw := data.([]interface{})
	if len(removeHeadersConfigRaw) == 1 {
		if v, ok := removeHeadersConfigRaw[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"configs": buildRemoveHeaderConfig(v["configs"]),
			}
			return params
		}
	}
	return nil
}

func buildRemoveHeaderConfig(data interface{}) []map[string]interface{} {
	removeHeaderConfigsRaw := data.(*schema.Set)
	if removeHeaderConfigsRaw.Len() == 0 {
		return nil
	}

	removeHeaderConfigs := make([]map[string]interface{}, 0, removeHeaderConfigsRaw.Len())
	for _, rawConfig := range removeHeaderConfigsRaw.List() {
		if v, ok := rawConfig.(map[string]interface{}); ok {
			removeHeaderConfigs = append(removeHeaderConfigs, map[string]interface{}{
				"key": v["key"].(string),
			})
		}
	}

	return removeHeaderConfigs
}

func buildTrafficLimitConfig(data interface{}) map[string]interface{} {
	trafficLimitConfigRaw := data.([]interface{})
	if len(trafficLimitConfigRaw) == 1 {
		if v, ok := trafficLimitConfigRaw[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"qps":               v["qps"].(int),
				"per_source_ip_qps": v["per_source_ip_qps"].(int),
				"burst":             v["burst"].(int),
			}
			return params
		}
	}

	return nil
}

func buildRedirectUrlConfig(d *schema.ResourceData) map[string]interface{} {
	if redirectUrlConfigRaw, ok := d.GetOk("redirect_url_config"); ok {
		if v, ok := redirectUrlConfigRaw.([]interface{})[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"protocol":              v["protocol"].(string),
				"host":                  v["host"].(string),
				"port":                  v["port"].(string),
				"path":                  v["path"].(string),
				"query":                 v["query"].(string),
				"status_code":           v["status_code"].(string),
				"insert_headers_config": buildInsertHeadersConfig(v["insert_headers_config"]),
				"remove_headers_config": buildRemoveHeadersConfig(v["remove_headers_config"]),
			}
			return params
		}
	}
	return nil
}

func buildFixedResponseConfig(d *schema.ResourceData) map[string]interface{} {
	if fixedResponseConfigRaw, ok := d.GetOk("fixed_response_config"); ok {
		if v, ok := fixedResponseConfigRaw.([]interface{})[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"status_code":           v["status_code"].(string),
				"content_type":          v["content_type"].(string),
				"message_body":          v["message_body"].(string),
				"insert_headers_config": buildInsertHeadersConfig(v["insert_headers_config"]),
				"remove_headers_config": buildRemoveHeadersConfig(v["remove_headers_config"]),
				"traffic_limit_config":  buildTrafficLimitConfig(v["traffic_limit_config"]),
			}
			return params
		}
	}
	return nil
}

func resourceL7PolicyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	l7Policy, err := getL7Policy(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB L7 policy")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("l7policy.name", l7Policy, nil)),
		d.Set("description", utils.PathSearch("l7policy.description", l7Policy, nil)),
		d.Set("action", utils.PathSearch("l7policy.action", l7Policy, nil)),
		d.Set("priority", utils.PathSearch("l7policy.priority", l7Policy, nil)),
		d.Set("listener_id", utils.PathSearch("l7policy.listener_id", l7Policy, nil)),
		d.Set("redirect_pool_id", utils.PathSearch("l7policy.redirect_pool_id", l7Policy, nil)),
		d.Set("redirect_listener_id", utils.PathSearch("l7policy.redirect_listener_id", l7Policy, nil)),
		d.Set("redirect_pools_config", flattenRedirectPoolsConfig(l7Policy)),
		d.Set("redirect_pools_sticky_session_config", flattenRedirectPoolsStickySessionConfig(l7Policy)),
		d.Set("redirect_pools_extend_config", flattenRedirectPoolsExtendConfig(l7Policy)),
		d.Set("redirect_url_config", flattenRedirectUrlConfig(l7Policy)),
		d.Set("fixed_response_config", flattenFixedResponseConfig(l7Policy)),
		d.Set("created_at", utils.PathSearch("l7policy.created_at", l7Policy, nil)),
		d.Set("updated_at", utils.PathSearch("l7policy.updated_at", l7Policy, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("l7policy.enterprise_project_id", l7Policy, nil)),
		d.Set("provisioning_status", utils.PathSearch("l7policy.provisioning_status", l7Policy, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRedirectPoolsConfig(l7policy interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("l7policy.redirect_pools_config", l7policy, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"pool_id": utils.PathSearch("pool_id", v, nil),
			"weight":  utils.PathSearch("weight", v, nil),
		})
	}
	return rst
}

func flattenRedirectPoolsStickySessionConfig(l7policy interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("l7policy.redirect_pools_sticky_session_config", l7policy, nil)
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

func flattenRedirectPoolsExtendConfig(l7policy interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("l7policy.redirect_pools_extend_config", l7policy, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"rewrite_url_enabled":   utils.PathSearch("rewrite_url_enable", curJson, nil),
			"rewrite_url_config":    flattenRewriteUrlConfig(curJson),
			"insert_headers_config": flattenInsertHeadersConfig(curJson),
			"remove_headers_config": flattenRemoveHeadersConfig(curJson),
			"traffic_limit_config":  flattenTrafficLimitConfig(curJson),
		},
	}
	return rst
}

func flattenRewriteUrlConfig(cfg interface{}) []map[string]interface{} {
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

func flattenInsertHeadersConfig(cfg interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("insert_headers_config", cfg, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"configs": flattenInsertHeaderConfigs(curJson),
		},
	}
	return rst
}

func flattenInsertHeaderConfigs(insertHeaderConfigs interface{}) []map[string]interface{} {
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

func flattenRemoveHeadersConfig(cfg interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("remove_headers_config", cfg, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"configs": flattenRemoveHeaderConfigs(curJson),
		},
	}
	return rst
}

func flattenRemoveHeaderConfigs(removeHeaderConfigs interface{}) []map[string]interface{} {
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

func flattenTrafficLimitConfig(cfg interface{}) []map[string]interface{} {
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

func flattenRedirectUrlConfig(l7policy interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("l7policy.redirect_url_config", l7policy, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"protocol":              utils.PathSearch("protocol", curJson, nil),
			"host":                  utils.PathSearch("host", curJson, nil),
			"port":                  utils.PathSearch("port", curJson, nil),
			"path":                  utils.PathSearch("path", curJson, nil),
			"query":                 utils.PathSearch("query", curJson, nil),
			"status_code":           utils.PathSearch("status_code", curJson, nil),
			"insert_headers_config": flattenInsertHeadersConfig(curJson),
			"remove_headers_config": flattenRemoveHeadersConfig(curJson),
		},
	}
	return rst
}

func flattenFixedResponseConfig(l7policy interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("l7policy.fixed_response_config", l7policy, nil)
	if curJson == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"status_code":           utils.PathSearch("status_code", curJson, nil),
			"content_type":          utils.PathSearch("content_type", curJson, nil),
			"message_body":          utils.PathSearch("message_body", curJson, nil),
			"insert_headers_config": flattenInsertHeadersConfig(curJson),
			"remove_headers_config": flattenRemoveHeadersConfig(curJson),
			"traffic_limit_config":  flattenTrafficLimitConfig(curJson),
		},
	}
	return rst
}

func resourceL7PolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/l7policies/{l7policy_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{l7policy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateL7PolicyBodyParams(d))
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ELB L7 policy: %s", err)
	}

	err = waitForL7Policy(ctx, client, d.Id(), "ACTIVE", nil, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7PolicyV3Read(ctx, d, meta)
}

func buildUpdateL7PolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                                 d.Get("name"),
		"description":                          d.Get("description"),
		"priority":                             d.Get("priority"),
		"redirect_pool_id":                     utils.ValueIgnoreEmpty(d.Get("redirect_pool_id")),
		"redirect_pools_config":                buildRedirectPoolsConfig(d),
		"redirect_pools_sticky_session_config": buildRedirectPoolsStickySessionConfig(d),
		"redirect_pools_extend_config":         buildRedirectPoolsExtendConfig(d),
		"redirect_listener_id":                 utils.ValueIgnoreEmpty(d.Get("redirect_listener_id")),
		"redirect_url_config":                  buildRedirectUrlConfig(d),
		"fixed_response_config":                buildFixedResponseConfig(d),
	}
	return map[string]interface{}{"l7policy": bodyParams}
}

func resourceL7PolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/l7policies/{l7policy_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{l7policy_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ELB L7 Policy")
	}

	err = waitForL7Policy(ctx, client, d.Id(), "DELETED", nil, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForL7Policy(ctx context.Context, client *golangsdk.ServiceClient, id string, target string,
	pending []string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{target},
		Pending:      pending,
		Refresh:      resourceL7PolicyRefreshFunc(client, id),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: L7 policy %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for L7 policy %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceL7PolicyRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := getL7Policy(client, id)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("l7policy.provisioning_status", policy, "")
		return policy, status.(string), nil
	}
}

func getL7Policy(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/l7policies/{l7policy_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{l7policy_id}", id)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}
