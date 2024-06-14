package elb

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/l7policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	action := d.Get("action").(string)
	createOpts := l7policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Action:      l7policies.Action(action),
		Priority:    int32(d.Get("priority").(int)),
		ListenerID:  d.Get("listener_id").(string),
	}
	if action == "REDIRECT_TO_POOL" {
		createOpts.RedirectPoolID = d.Get("redirect_pool_id").(string)
		createOpts.RedirectPoolsConfig = buildRedirectPoolsConfig(d)
		createOpts.RedirectPoolsStickySessionConfig = buildRedirectPoolsStickySessionConfig(d)
		createOpts.RedirectPoolsExtendConfig = buildRedirectPoolsExtendConfig(d)
	} else if action == "REDIRECT_TO_LISTENER" {
		createOpts.RedirectListenerID = d.Get("redirect_listener_id").(string)
	} else if action == "REDIRECT_TO_URL" {
		createOpts.RedirectUrlConfig = buildRedirectUrlConfig(d)
	} else {
		createOpts.FixedResponseConfig = buildFixedResponseConfig(d)
	}

	l7Policy, err := l7policies.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating L7 Policy: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	// Wait for L7 Policy to become active before continuing
	err = waitForElbV3Policy(ctx, elbClient, l7Policy.ID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Policy.ID)

	return resourceL7PolicyV3Read(ctx, d, meta)
}

func buildRedirectPoolsConfig(d *schema.ResourceData) []*l7policies.RedirectPoolsConfig {
	var redirectPoolsConfig []*l7policies.RedirectPoolsConfig
	redirectPoolsConfigRaw := d.Get("redirect_pools_config").(*schema.Set).List()
	for _, redirectPoolConfigRaw := range redirectPoolsConfigRaw {
		v := redirectPoolConfigRaw.(map[string]interface{})
		redirectPoolsConfig = append(redirectPoolsConfig, &l7policies.RedirectPoolsConfig{
			PoolId: v["pool_id"].(string),
			Weight: v["weight"].(int),
		})
	}
	return redirectPoolsConfig
}

func buildRedirectPoolsStickySessionConfig(d *schema.ResourceData) *l7policies.RedirectPoolsStickySessionConfig {
	var redirectPoolsStickySessionConfig *l7policies.RedirectPoolsStickySessionConfig
	redirectPoolsStickySessionConfigRaw := d.Get("redirect_pools_sticky_session_config").([]interface{})
	if len(redirectPoolsStickySessionConfigRaw) == 1 {
		if v, ok := redirectPoolsStickySessionConfigRaw[0].(map[string]interface{}); ok {
			redirectPoolsStickySessionConfig = &l7policies.RedirectPoolsStickySessionConfig{
				Enable:  v["enable"].(bool),
				Timeout: v["timeout"].(int),
			}
		}
	}
	return redirectPoolsStickySessionConfig
}

func buildRedirectPoolsExtendConfig(d *schema.ResourceData) *l7policies.RedirectPoolsExtendConfig {
	var redirectPoolsExtendConfig *l7policies.RedirectPoolsExtendConfig
	redirectPoolsExtendConfigRaw := d.Get("redirect_pools_extend_config").([]interface{})
	if len(redirectPoolsExtendConfigRaw) == 1 {
		if v, ok := redirectPoolsExtendConfigRaw[0].(map[string]interface{}); ok {
			redirectPoolsExtendConfig = &l7policies.RedirectPoolsExtendConfig{
				RewriteUrlEnable:    v["rewrite_url_enabled"].(bool),
				RewriteUrlConfig:    buildRewriteUrlConfig(v["rewrite_url_config"]),
				InsertHeadersConfig: buildInsertHeadersConfig(v["insert_headers_config"]),
				RemoveHeadersConfig: buildRemoveHeadersConfig(v["remove_headers_config"]),
				TrafficLimitConfig:  buildTrafficLimitConfig(v["traffic_limit_config"]),
			}
		}
	}
	return redirectPoolsExtendConfig
}

func buildRewriteUrlConfig(data interface{}) *l7policies.RewriteUrlConfig {
	var rewriteUrlConfig *l7policies.RewriteUrlConfig
	rewriteUrlConfigRaw := data.([]interface{})
	if len(rewriteUrlConfigRaw) == 1 {
		if v, ok := rewriteUrlConfigRaw[0].(map[string]interface{}); ok {
			rewriteUrlConfig = &l7policies.RewriteUrlConfig{
				Host:  v["host"].(string),
				Path:  v["path"].(string),
				Query: v["query"].(string),
			}
		}
	}
	return rewriteUrlConfig
}

func buildInsertHeadersConfig(data interface{}) *l7policies.InsertHeadersConfig {
	var insertHeadersConfig *l7policies.InsertHeadersConfig
	insertHeadersConfigRaw := data.([]interface{})
	if len(insertHeadersConfigRaw) == 1 {
		if v, ok := insertHeadersConfigRaw[0].(map[string]interface{}); ok {
			insertHeadersConfig = &l7policies.InsertHeadersConfig{
				Configs: buildInsertHeaderConfig(v["configs"]),
			}
		}
	}
	return insertHeadersConfig
}

func buildInsertHeaderConfig(data interface{}) []*l7policies.InsertHeaderConfig {
	var insertHeaderConfigs []*l7policies.InsertHeaderConfig
	insertHeaderConfigsRaw := data.(*schema.Set).List()
	for _, insertHeaderConfigRaw := range insertHeaderConfigsRaw {
		v := insertHeaderConfigRaw.(map[string]interface{})
		insertHeaderConfigs = append(insertHeaderConfigs, &l7policies.InsertHeaderConfig{
			Key:       v["key"].(string),
			ValueType: v["value_type"].(string),
			Value:     v["value"].(string),
		})
	}
	return insertHeaderConfigs
}

func buildRemoveHeadersConfig(data interface{}) *l7policies.RemoveHeadersConfig {
	var removeHeadersConfig *l7policies.RemoveHeadersConfig
	removeHeadersConfigRaw := data.([]interface{})
	if len(removeHeadersConfigRaw) == 1 {
		if v, ok := removeHeadersConfigRaw[0].(map[string]interface{}); ok {
			removeHeadersConfig = &l7policies.RemoveHeadersConfig{
				Configs: buildRemoveHeaderConfig(v["configs"]),
			}
		}
	}
	return removeHeadersConfig
}

func buildRemoveHeaderConfig(data interface{}) []*l7policies.RemoveHeaderConfig {
	var removeHeaderConfigs []*l7policies.RemoveHeaderConfig
	removeHeaderConfigsRaw := data.(*schema.Set).List()
	for _, removeHeaderConfigRaw := range removeHeaderConfigsRaw {
		v := removeHeaderConfigRaw.(map[string]interface{})
		removeHeaderConfigs = append(removeHeaderConfigs, &l7policies.RemoveHeaderConfig{
			Key: v["key"].(string),
		})
	}
	return removeHeaderConfigs
}

func buildTrafficLimitConfig(data interface{}) *l7policies.TrafficLimitConfig {
	var trafficLimitConfig *l7policies.TrafficLimitConfig
	trafficLimitConfigRaw := data.([]interface{})
	if len(trafficLimitConfigRaw) == 1 {
		if v, ok := trafficLimitConfigRaw[0].(map[string]interface{}); ok {
			trafficLimitConfig = &l7policies.TrafficLimitConfig{
				Qps:            v["qps"].(int),
				PerSourceIpQps: v["per_source_ip_qps"].(int),
				Burst:          v["burst"].(int),
			}
		}
	}
	return trafficLimitConfig
}

func buildRedirectUrlConfig(d *schema.ResourceData) *l7policies.RedirectUrlConfig {
	var redirectUrlConfig *l7policies.RedirectUrlConfig
	redirectUrlConfigRaw := d.Get("redirect_url_config").([]interface{})
	if len(redirectUrlConfigRaw) == 1 {
		if v, ok := redirectUrlConfigRaw[0].(map[string]interface{}); ok {
			redirectUrlConfig = &l7policies.RedirectUrlConfig{
				Protocol:            v["protocol"].(string),
				Host:                v["host"].(string),
				Port:                v["port"].(string),
				Path:                v["path"].(string),
				Query:               v["query"].(string),
				StatusCode:          v["status_code"].(string),
				InsertHeadersConfig: buildInsertHeadersConfig(v["insert_headers_config"]),
				RemoveHeadersConfig: buildRemoveHeadersConfig(v["remove_headers_config"]),
			}
		}
	}
	return redirectUrlConfig
}

func buildFixedResponseConfig(d *schema.ResourceData) *l7policies.FixedResponseConfig {
	var fixedResponseConfig *l7policies.FixedResponseConfig
	fixedResponseConfigRaw := d.Get("fixed_response_config").([]interface{})
	if len(fixedResponseConfigRaw) == 1 {
		if v, ok := fixedResponseConfigRaw[0].(map[string]interface{}); ok {
			fixedResponseConfig = &l7policies.FixedResponseConfig{
				StatusCode:          v["status_code"].(string),
				ContentType:         v["content_type"].(string),
				MessageBody:         v["message_body"].(string),
				InsertHeadersConfig: buildInsertHeadersConfig(v["insert_headers_config"]),
				RemoveHeadersConfig: buildRemoveHeadersConfig(v["remove_headers_config"]),
				TrafficLimitConfig:  buildTrafficLimitConfig(v["traffic_limit_config"]),
			}
		}
	}
	return fixedResponseConfig
}

func resourceL7PolicyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	l7Policy, err := l7policies.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "L7 Policy")
	}

	mErr := multierror.Append(nil,
		d.Set("description", l7Policy.Description),
		d.Set("name", l7Policy.Name),
		d.Set("action", l7Policy.Action),
		d.Set("priority", l7Policy.Priority),
		d.Set("listener_id", l7Policy.ListenerID),
		d.Set("redirect_pool_id", l7Policy.RedirectPoolID),
		d.Set("redirect_listener_id", l7Policy.RedirectListenerID),
		d.Set("redirect_pools_config", flattenRedirectPoolsConfig(l7Policy)),
		d.Set("redirect_pools_sticky_session_config", flattenRedirectPoolsStickySessionConfig(l7Policy)),
		d.Set("redirect_pools_extend_config", flattenRedirectPoolsExtendConfig(l7Policy)),
		d.Set("redirect_url_config", flattenRedirectUrlConfig(l7Policy)),
		d.Set("fixed_response_config", flattenFixedResponseConfig(l7Policy)),
		d.Set("region", cfg.GetRegion(d)),
		d.Set("created_at", l7Policy.CreatedAt),
		d.Set("updated_at", l7Policy.UpdatedAt),
		d.Set("provisioning_status", l7Policy.ProvisioningStatus),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB l7policy fields: %s", err)
	}

	return nil
}

func flattenRedirectPoolsConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var redirectPoolsConfig []map[string]interface{}
	if l7policy.RedirectPoolsConfig != nil {
		redirectPoolsConfig = make([]map[string]interface{}, 0, len(l7policy.RedirectPoolsConfig))
		for _, redirectPoolConfig := range l7policy.RedirectPoolsConfig {
			redirectPoolsConfig = append(redirectPoolsConfig, map[string]interface{}{
				"pool_id": redirectPoolConfig.PoolId,
				"weight":  redirectPoolConfig.Weight,
			})
		}
	}
	return redirectPoolsConfig
}

func flattenRedirectPoolsStickySessionConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var redirectPoolsStickySessionConfig []map[string]interface{}
	if l7policy.RedirectPoolsStickySessionConfig != nil {
		redirectPoolsStickySessionConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["enable"] = l7policy.RedirectPoolsStickySessionConfig.Enable
		params["timeout"] = l7policy.RedirectPoolsStickySessionConfig.Timeout
		redirectPoolsStickySessionConfig[0] = params
	}
	return redirectPoolsStickySessionConfig
}

func flattenRedirectPoolsExtendConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var redirectPoolsExtendConfig []map[string]interface{}
	if l7policy.RedirectPoolsExtendConfig != nil {
		redirectPoolsExtendConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["rewrite_url_enabled"] = l7policy.RedirectPoolsExtendConfig.RewriteUrlEnable
		params["rewrite_url_config"] = flattenRewriteUrlConfig(l7policy)
		params["insert_headers_config"] = flattenInsertHeadersConfig(l7policy.RedirectPoolsExtendConfig.InsertHeadersConfig)
		params["remove_headers_config"] = flattenRemoveHeadersConfig(l7policy.RedirectPoolsExtendConfig.RemoveHeadersConfig)
		params["traffic_limit_config"] = flattenTrafficLimitConfig(l7policy.RedirectPoolsExtendConfig.TrafficLimitConfig)
		redirectPoolsExtendConfig[0] = params
	}
	return redirectPoolsExtendConfig
}

func flattenRewriteUrlConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var rewriteUrlConfig []map[string]interface{}
	if l7policy.RedirectPoolsExtendConfig.RewriteUrlConfig != nil {
		rewriteUrlConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["host"] = l7policy.RedirectPoolsExtendConfig.RewriteUrlConfig.Host
		params["path"] = l7policy.RedirectPoolsExtendConfig.RewriteUrlConfig.Path
		params["query"] = l7policy.RedirectPoolsExtendConfig.RewriteUrlConfig.Query
		rewriteUrlConfig[0] = params
	}
	return rewriteUrlConfig
}

func flattenInsertHeadersConfig(cfg *l7policies.InsertHeadersConfig) []map[string]interface{} {
	var insertHeadersConfig []map[string]interface{}
	if cfg != nil {
		insertHeadersConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["configs"] = flattenInsertHeaderConfigs(cfg.Configs)
		insertHeadersConfig[0] = params
	}
	return insertHeadersConfig
}

func flattenInsertHeaderConfigs(insertHeaderConfigs []*l7policies.InsertHeaderConfig) []map[string]interface{} {
	var configs []map[string]interface{}
	if len(insertHeaderConfigs) > 0 {
		configs = make([]map[string]interface{}, 0, len(insertHeaderConfigs))
		for _, v := range insertHeaderConfigs {
			configs = append(configs, map[string]interface{}{
				"key":        v.Key,
				"value_type": v.ValueType,
				"value":      v.Value,
			})
		}
	}
	return configs
}

func flattenRemoveHeadersConfig(cfg *l7policies.RemoveHeadersConfig) []map[string]interface{} {
	var removeHeadersConfig []map[string]interface{}
	if cfg != nil {
		removeHeadersConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["configs"] = flattenRemoveHeaderConfigs(cfg.Configs)
		removeHeadersConfig[0] = params
	}
	return removeHeadersConfig
}

func flattenRemoveHeaderConfigs(removeHeaderConfigs []*l7policies.RemoveHeaderConfig) []map[string]interface{} {
	var configs []map[string]interface{}
	if len(removeHeaderConfigs) > 0 {
		configs = make([]map[string]interface{}, 0, len(removeHeaderConfigs))
		for _, v := range removeHeaderConfigs {
			configs = append(configs, map[string]interface{}{
				"key": v.Key,
			})
		}
	}
	return configs
}

func flattenTrafficLimitConfig(cfg *l7policies.TrafficLimitConfig) []map[string]interface{} {
	var trafficLimitConfig []map[string]interface{}
	if cfg != nil {
		trafficLimitConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["qps"] = cfg.Qps
		params["per_source_ip_qps"] = cfg.PerSourceIpQps
		params["burst"] = cfg.Burst
		trafficLimitConfig[0] = params
	}
	return trafficLimitConfig
}

func flattenRedirectUrlConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var redirectUrlConfig []map[string]interface{}
	if l7policy.RedirectUrlConfig != nil {
		redirectUrlConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["protocol"] = l7policy.RedirectUrlConfig.Protocol
		params["host"] = l7policy.RedirectUrlConfig.Host
		params["port"] = l7policy.RedirectUrlConfig.Port
		params["path"] = l7policy.RedirectUrlConfig.Path
		params["query"] = l7policy.RedirectUrlConfig.Query
		params["status_code"] = l7policy.RedirectUrlConfig.StatusCode
		params["insert_headers_config"] = flattenInsertHeadersConfig(l7policy.RedirectUrlConfig.InsertHeadersConfig)
		params["remove_headers_config"] = flattenRemoveHeadersConfig(l7policy.RedirectUrlConfig.RemoveHeadersConfig)
		redirectUrlConfig[0] = params
	}
	return redirectUrlConfig
}

func flattenFixedResponseConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var fixedResponseConfig []map[string]interface{}
	if l7policy.FixedResponseConfig != nil {
		fixedResponseConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["status_code"] = l7policy.FixedResponseConfig.StatusCode
		params["content_type"] = l7policy.FixedResponseConfig.ContentType
		params["message_body"] = l7policy.FixedResponseConfig.MessageBody
		params["insert_headers_config"] = flattenInsertHeadersConfig(l7policy.FixedResponseConfig.InsertHeadersConfig)
		params["remove_headers_config"] = flattenRemoveHeadersConfig(l7policy.FixedResponseConfig.RemoveHeadersConfig)
		params["traffic_limit_config"] = flattenTrafficLimitConfig(l7policy.FixedResponseConfig.TrafficLimitConfig)
		fixedResponseConfig[0] = params
	}
	return fixedResponseConfig
}

func resourceL7PolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	var updateOpts l7policies.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("priority") {
		priority := d.Get("priority").(int)
		updateOpts.Priority = int32(priority)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("redirect_pool_id") {
		redirectPoolID := d.Get("redirect_pool_id").(string)
		updateOpts.RedirectPoolID = &redirectPoolID
	}
	if d.HasChange("redirect_pools_config") {
		updateOpts.RedirectPoolsConfig = buildRedirectPoolsConfig(d)
	}
	if d.HasChange("redirect_pools_sticky_session_config") {
		updateOpts.RedirectPoolsStickySessionConfig = buildRedirectPoolsStickySessionConfig(d)
	}
	if d.HasChange("redirect_pools_extend_config") {
		updateOpts.RedirectPoolsExtendConfig = buildRedirectPoolsExtendConfig(d)
	}
	if d.HasChange("redirect_listener_id") {
		redirectListenerID := d.Get("redirect_listener_id").(string)
		updateOpts.RedirectListenerID = &redirectListenerID
	}
	if d.HasChange("redirect_url_config") {
		updateOpts.RedirectUrlConfig = buildRedirectUrlConfig(d)
	}
	if d.HasChange("fixed_response_config") {
		updateOpts.FixedResponseConfig = buildFixedResponseConfig(d)
	}

	_, err = l7policies.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update L7 Policy %s: %s", d.Id(), err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForElbV3Policy(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7PolicyV3Read(ctx, d, meta)
}

func resourceL7PolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	err = l7policies.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting L7 Policy")
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbV3Policy(ctx, elbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForElbV3Policy(ctx context.Context, elbClient *golangsdk.ServiceClient, id string, target string,
	pending []string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{target},
		Pending:      pending,
		Refresh:      resourceElbV3PolicyRefreshFunc(elbClient, id),
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
				return fmt.Errorf("error: policy %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for policy %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceElbV3PolicyRefreshFunc(elbClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := l7policies.Get(elbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return policy, policy.ProvisioningStatus, nil
	}
}
