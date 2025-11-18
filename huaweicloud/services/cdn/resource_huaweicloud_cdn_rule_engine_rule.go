package cdn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ruleEngineRuleNonUpdatableParams = []string{"domain_name"}

// @API CDN POST /v1.0/cdn/configuration/domains/{domain_name}/rules
// @API CDN GET /v1.0/cdn/configuration/domains/{domain_name}/rules
// @API CDN PUT /v1.0/cdn/configuration/domains/{domain_name}/rules/{rule_id}
// @API CDN DELETE /v1.0/cdn/configuration/domains/{domain_name}/rules/{rule_id}
func ResourceRuleEngineRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleEngineRuleCreate,
		ReadContext:   resourceRuleEngineRuleRead,
		UpdateContext: resourceRuleEngineRuleUpdate,
		DeleteContext: resourceRuleEngineRuleDelete,

		CustomizeDiff: config.FlexibleForceNew(ruleEngineRuleNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceRuleEngineRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the rule engine rule is located.`,
			},

			// Required parameters.
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The accelerated domain name to which the rule engine rule belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the rule engine rule.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Whether to enable the rule engine rule.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The priority of the rule engine rule.`,
			},

			// Optional parameters.
			"conditions": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The trigger conditions of the rule engine rule, in JSON format.`,
			},
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flexible_origin": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `The list of flexible origin configurations.`,
							Elem:        ruleEngineRuleActionsFlexibleOriginSchema(),
						},
						"origin_request_header": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `The list of origin request header configurations.`,
							Elem:        ruleEngineRuleActionsOriginRequestHeaderSchema(),
						},
						"http_response_header": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `The list of HTTP response header configurations.`,
							Elem:        ruleEngineRuleActionsHttpResponseHeaderSchema(),
						},
						"access_control": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `The access control configuration.`,
							Elem:        ruleEngineRuleActionsAccessControlSchema(),
						},
						"request_limit_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `The request rate limit configuration.`,
							Elem:        ruleEngineRuleActionsRequestLimitRuleSchema(),
						},
						"origin_request_url_rewrite": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `The origin request URL rewrite configuration.`,
							Elem:        ruleEngineRuleActionsOriginRequestUrlRewriteSchema(),
						},
						"cache_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `The cache rule configuration.`,
							Elem:        ruleEngineRuleActionsCacheRuleSchema(),
						},
						"request_url_rewrite": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `The access URL rewrite configuration.`,
							Elem:        ruleEngineRuleActionsRequestUrlRewriteSchema(),
						},
						"browser_cache_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `The browser cache rule configuration.`,
							Elem:        ruleEngineRuleActionsBrowserCacheRuleSchema(),
						},
						"error_code_cache": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `The list of error code cache configurations.`,
							Elem:        ruleEngineRuleActionsErrorCodeCacheSchema(),
						},
						"origin_range": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `The origin range configuration.`,
							Elem:        ruleEngineRuleActionsOriginRangeSchema(),
						},
					},
				},
				Description: `The list of actions to be performed when the rule engine rule is met.`,
			},

			// Internal parameter.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},

			// Internal attribute.
			"conditions_origin": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'conditions'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func ruleEngineRuleActionsFlexibleOriginSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"sources_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The source type.`,
			},
			"ip_or_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The origin IP or domain name.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The origin priority.`,
			},
			"weight": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The origin weight.`,
			},

			// Optional parameters.
			"obs_bucket_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The OBS bucket type.`,
			},
			"bucket_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The third-party object storage access key.`,
			},
			"bucket_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The third-party object storage secret key.`,
			},
			"bucket_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The third-party object storage region.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The third-party object storage name.`,
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The origin host name.`,
			},
			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The origin protocol.`,
			},
			"http_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The HTTP port number.`,
			},
			"https_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The HTTPS port number.`,
			},
		},
	}
}

func ruleEngineRuleActionsOriginRequestHeaderSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The back-to-origin request header parameter name.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The back-to-origin request header setting type.`,
			},

			// Optional parameters.
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The back-to-origin request header parameter value.`,
			},
		},
	}
}

func ruleEngineRuleActionsHttpResponseHeaderSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The HTTP response header parameter name.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The operation type of setting HTTP response header.`,
			},

			// Optional parameters.
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The HTTP response header parameter value.`,
			},
		},
	}
}

func ruleEngineRuleActionsAccessControlSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The access control type.`,
			},
		},
	}
}

func ruleEngineRuleActionsRequestLimitRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"limit_rate_after": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The rate limit condition.`,
			},
			"limit_rate_value": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The rate limit value.`,
			},
		},
	}
}

func ruleEngineRuleActionsOriginRequestUrlRewriteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"rewrite_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The rewrite type.`,
			},
			"target_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The target URL.`,
			},

			// Optional parameters.
			"source_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The source URL to be rewritten.`,
			},
		},
	}
}

func ruleEngineRuleActionsCacheRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The cache expiration time.`,
			},
			"ttl_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The cache expiration time unit.`,
			},
			"follow_origin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The cache expiration time source.`,
			},

			// Optional parameters.
			"force_cache": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable forced caching.`,
			},
		},
	}
}

func ruleEngineRuleActionsRequestUrlRewriteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"redirect_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The redirect URL.`,
			},
			"execution_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The execution mode.`,
			},

			// Optional parameters.
			"redirect_status_code": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The redirect status code.`,
			},
			"redirect_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The redirect host.`,
			},
		},
	}
}

func ruleEngineRuleActionsBrowserCacheRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"cache_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The cache effective type.`,
			},

			// Optional parameters.
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The cache expiration time.`,
			},
			"ttl_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The cache expiration time unit.`,
			},
		},
	}
}

func ruleEngineRuleActionsErrorCodeCacheSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"code": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The error code to be cached.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The error code cache time.`,
			},
		},
	}
}

func ruleEngineRuleActionsOriginRangeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The origin range status.`,
			},
		},
	}
}

func buildRuleEngineRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":       d.Get("name").(string),
		"status":     d.Get("status").(string),
		"priority":   d.Get("priority").(int),
		"conditions": utils.StringToJson(d.Get("conditions").(string)),
		"actions":    buildRuleEngineRuleActionsBodyParams(d.Get("actions").([]interface{})),
	}
}

func buildRuleEngineRuleActionsBodyParams(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"flexible_origin": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsFlexibleOriginBodyParams(
				utils.PathSearch("flexible_origin", item, make([]interface{}, 0)).([]interface{}))),
			"origin_request_header": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsOriginRequestHeaderBodyParams(
				utils.PathSearch("origin_request_header", item, make([]interface{}, 0)).([]interface{}))),
			"http_response_header": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsHttpResponseHeaderBodyParams(
				utils.PathSearch("http_response_header", item, make([]interface{}, 0)).([]interface{}))),
			"access_control": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsAccessControlBodyParams(
				utils.PathSearch("access_control", item, make([]interface{}, 0)).([]interface{}))),
			"request_limit_rule": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsRequestLimitRuleBodyParams(
				utils.PathSearch("request_limit_rule", item, make([]interface{}, 0)).([]interface{}))),
			"origin_request_url_rewrite": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsOriginRequestUrlRewriteBodyParams(
				utils.PathSearch("origin_request_url_rewrite", item, make([]interface{}, 0)).([]interface{}))),
			"cache_rule": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsCacheRuleBodyParams(
				utils.PathSearch("cache_rule", item, make([]interface{}, 0)).([]interface{}))),
			"request_url_rewrite": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsRequestUrlRewriteBodyParams(
				utils.PathSearch("request_url_rewrite", item, make([]interface{}, 0)).([]interface{}))),
			"browser_cache_rule": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsBrowserCacheRuleBodyParams(
				utils.PathSearch("browser_cache_rule", item, make([]interface{}, 0)).([]interface{}))),
			"error_code_cache": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsErrorCodeCacheBodyParams(
				utils.PathSearch("error_code_cache", item, make([]interface{}, 0)).([]interface{}))),
			"origin_range": utils.ValueIgnoreEmpty(buildRuleEngineRuleActionsOriginRangeBodyParams(
				utils.PathSearch("origin_range", item, make([]interface{}, 0)).([]interface{}))),
		}))
	}

	return result
}

func buildRuleEngineRuleActionsFlexibleOriginBodyParams(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"sources_type":      utils.ValueIgnoreEmpty(utils.PathSearch("sources_type", item, nil)),
			"ip_or_domain":      utils.ValueIgnoreEmpty(utils.PathSearch("ip_or_domain", item, nil)),
			"priority":          utils.ValueIgnoreEmpty(utils.PathSearch("priority", item, nil)),
			"weight":            utils.ValueIgnoreEmpty(utils.PathSearch("weight", item, nil)),
			"obs_bucket_type":   utils.ValueIgnoreEmpty(utils.PathSearch("obs_bucket_type", item, nil)),
			"bucket_access_key": utils.ValueIgnoreEmpty(utils.PathSearch("bucket_access_key", item, nil)),
			"bucket_secret_key": utils.ValueIgnoreEmpty(utils.PathSearch("bucket_secret_key", item, nil)),
			"bucket_region":     utils.ValueIgnoreEmpty(utils.PathSearch("bucket_region", item, nil)),
			"bucket_name":       utils.ValueIgnoreEmpty(utils.PathSearch("bucket_name", item, nil)),
			"host_name":         utils.ValueIgnoreEmpty(utils.PathSearch("host_name", item, nil)),
			"origin_protocol":   utils.ValueIgnoreEmpty(utils.PathSearch("origin_protocol", item, nil)),
			"http_port":         utils.ValueIgnoreEmpty(utils.PathSearch("http_port", item, nil)),
			"https_port":        utils.ValueIgnoreEmpty(utils.PathSearch("https_port", item, nil)),
		}))
	}

	return result
}

func buildRuleEngineRuleActionsOriginRequestHeaderBodyParams(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"action": utils.ValueIgnoreEmpty(utils.PathSearch("action", item, nil)),
			"name":   utils.ValueIgnoreEmpty(utils.PathSearch("name", item, nil)),
			"value":  utils.ValueIgnoreEmpty(utils.PathSearch("value", item, nil)),
		}))
	}

	return result
}

func buildRuleEngineRuleActionsHttpResponseHeaderBodyParams(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"action": utils.ValueIgnoreEmpty(utils.PathSearch("action", item, nil)),
			"name":   utils.ValueIgnoreEmpty(utils.PathSearch("name", item, nil)),
			"value":  utils.ValueIgnoreEmpty(utils.PathSearch("value", item, nil)),
		}))
	}

	return result
}

func buildRuleEngineRuleActionsAccessControlBodyParams(items []interface{}) map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"type": utils.ValueIgnoreEmpty(utils.PathSearch("type", items[0], nil)),
	})
}

func buildRuleEngineRuleActionsRequestLimitRuleBodyParams(items []interface{}) map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"limit_rate_after": utils.ValueIgnoreEmpty(utils.PathSearch("limit_rate_after", items[0], nil)),
		"limit_rate_value": utils.ValueIgnoreEmpty(utils.PathSearch("limit_rate_value", items[0], nil)),
	})
}

func buildRuleEngineRuleActionsOriginRequestUrlRewriteBodyParams(items []interface{}) map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"rewrite_type": utils.ValueIgnoreEmpty(utils.PathSearch("rewrite_type", items[0], nil)),
		"target_url":   utils.ValueIgnoreEmpty(utils.PathSearch("target_url", items[0], nil)),
		"source_url":   utils.ValueIgnoreEmpty(utils.PathSearch("source_url", items[0], nil)),
	})
}

func buildRuleEngineRuleActionsCacheRuleBodyParams(items []interface{}) map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"ttl":           utils.ValueIgnoreEmpty(utils.PathSearch("ttl", items[0], nil)),
		"ttl_unit":      utils.ValueIgnoreEmpty(utils.PathSearch("ttl_unit", items[0], nil)),
		"follow_origin": utils.ValueIgnoreEmpty(utils.PathSearch("follow_origin", items[0], nil)),
		"force_cache":   utils.ValueIgnoreEmpty(utils.PathSearch("force_cache", items[0], nil)),
	})
}

func buildRuleEngineRuleActionsRequestUrlRewriteBodyParams(items []interface{}) map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	return map[string]interface{}{
		"redirect_url":         utils.ValueIgnoreEmpty(utils.PathSearch("redirect_url", items[0], nil)),
		"execution_mode":       utils.ValueIgnoreEmpty(utils.PathSearch("execution_mode", items[0], nil)),
		"redirect_status_code": utils.ValueIgnoreEmpty(utils.PathSearch("redirect_status_code", items[0], nil)),
		"redirect_host":        utils.ValueIgnoreEmpty(utils.PathSearch("redirect_host", items[0], nil)),
	}
}

func buildRuleEngineRuleActionsBrowserCacheRuleBodyParams(items []interface{}) map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"cache_type": utils.ValueIgnoreEmpty(utils.PathSearch("cache_type", items[0], nil)),
		"ttl":        utils.ValueIgnoreEmpty(utils.PathSearch("ttl", items[0], nil)),
		"ttl_unit":   utils.ValueIgnoreEmpty(utils.PathSearch("ttl_unit", items[0], nil)),
	})
}

func buildRuleEngineRuleActionsErrorCodeCacheBodyParams(items []interface{}) []interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"code": utils.ValueIgnoreEmpty(utils.PathSearch("code", item, nil)),
			"ttl":  utils.ValueIgnoreEmpty(utils.PathSearch("ttl", item, nil)),
		})
	}

	return result
}

func buildRuleEngineRuleActionsOriginRangeBodyParams(items []interface{}) map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"status": utils.ValueIgnoreEmpty(utils.PathSearch("status", items[0], nil)),
	})
}

func createRuleEngineRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1.0/cdn/configuration/domains/{domain_name}/rules"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_name}", d.Get("domain_name").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRuleEngineRuleBodyParams(d)),
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func listRuleEngineRules(client *golangsdk.ServiceClient, domainName string) ([]interface{}, error) {
	httpUrl := "v1.0/cdn/configuration/domains/{domain_name}/rules"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_name}", domainName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	rules := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
	return rules, nil
}

func GetRuleEngineRuleByName(client *golangsdk.ServiceClient, domainName string, ruleName string) (interface{}, error) {
	rules, err := listRuleEngineRules(client, domainName)
	if err != nil {
		return nil, err
	}

	rule := utils.PathSearch(fmt.Sprintf("[?name =='%s']|[0]", ruleName), rules, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/cdn/configuration/domains/{domain_name}/rules",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the rule with name '%s' has been removed", ruleName)),
			},
		}
	}
	return rule, nil
}

func resourceRuleEngineRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		domainName = d.Get("domain_name").(string)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	if err := createRuleEngineRule(client, d); err != nil {
		return diag.Errorf("error creating CDN rule engine rule: %s", err)
	}

	ruleName := d.Get("name").(string)
	rule, err := GetRuleEngineRuleByName(client, domainName, ruleName)
	if err != nil {
		return diag.Errorf("unable to find the created rule with name '%s': %s", ruleName, err)
	}

	ruleId := utils.PathSearch("rule_id", rule, "").(string)
	if ruleId == "" {
		return diag.Errorf("unable to find the rule ID from the API response")
	}

	d.SetId(ruleId)

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, []string{"conditions"})
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceRuleEngineRuleRead(ctx, d, meta)
}

func flattenRuleEngineRuleActionsFlexibleOrigin(items, originItems []interface{}) []interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
	for i, item := range items {
		result = append(result, map[string]interface{}{
			"priority":          utils.PathSearch("priority", item, nil),
			"weight":            utils.PathSearch("weight", item, nil),
			"sources_type":      utils.PathSearch("sources_type", item, nil),
			"ip_or_domain":      utils.PathSearch("ip_or_domain", item, nil),
			"obs_bucket_type":   utils.PathSearch("obs_bucket_type", item, nil),
			"bucket_access_key": utils.PathSearch("bucket_access_key", item, nil),
			"bucket_secret_key": utils.PathSearch(fmt.Sprintf("[%d].bucket_secret_key", i), originItems, nil),
			"bucket_region":     utils.PathSearch("bucket_region", item, nil),
			"bucket_name":       utils.PathSearch("bucket_name", item, nil),
			"host_name":         utils.PathSearch("host_name", item, nil),
			"origin_protocol":   utils.PathSearch("origin_protocol", item, nil),
			"http_port":         utils.PathSearch("http_port", item, nil),
			"https_port":        utils.PathSearch("https_port", item, nil),
		})
	}

	return result
}

func flattenRuleEngineRuleActionsOriginRequestHeader(items []interface{}) []interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"action": utils.PathSearch("action", item, nil),
			"name":   utils.PathSearch("name", item, nil),
			"value":  utils.PathSearch("value", item, nil),
		})
	}

	return result
}

func flattenRuleEngineRuleActionsHttpResponseHeader(items []interface{}) []interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"action": utils.PathSearch("action", item, nil),
			"name":   utils.PathSearch("name", item, nil),
			"value":  utils.PathSearch("value", item, nil),
		})
	}

	return result
}

func flattenRuleEngineRuleActionsAccessControl(item map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"type": utils.PathSearch("type", item, nil),
		},
	}
}

func flattenRuleEngineRuleActionsRequestLimitRule(item map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"limit_rate_after": utils.PathSearch("limit_rate_after", item, nil),
			"limit_rate_value": utils.PathSearch("limit_rate_value", item, nil),
		},
	}
}

func flattenRuleEngineRuleActionsOriginRequestUrlRewrite(item map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"rewrite_type": utils.PathSearch("rewrite_type", item, nil),
			"source_url":   utils.PathSearch("source_url", item, nil),
			"target_url":   utils.PathSearch("target_url", item, nil),
		},
	}
}

func flattenRuleEngineRuleActionsCacheRule(item map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"ttl":           utils.PathSearch("ttl", item, nil),
			"ttl_unit":      utils.PathSearch("ttl_unit", item, nil),
			"follow_origin": utils.PathSearch("follow_origin", item, nil),
			"force_cache":   utils.PathSearch("force_cache", item, nil),
		},
	}
}

func flattenRuleEngineRuleActionsRequestUrlRewrite(item map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"redirect_url":         utils.PathSearch("redirect_url", item, nil),
			"execution_mode":       utils.PathSearch("execution_mode", item, nil),
			"redirect_status_code": utils.PathSearch("redirect_status_code", item, nil),
			"redirect_host":        utils.PathSearch("redirect_host", item, nil),
		},
	}
}

func flattenRuleEngineRuleActionsBrowserCacheRule(item map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"cache_type": utils.PathSearch("cache_type", item, nil),
			"ttl":        utils.PathSearch("ttl", item, nil),
			"ttl_unit":   utils.PathSearch("ttl_unit", item, nil),
		},
	}
}

func flattenRuleEngineRuleActionsErrorCodeCache(items []interface{}) []interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"code": utils.PathSearch("code", item, nil),
			"ttl":  utils.PathSearch("ttl", item, nil),
		})
	}

	return result
}

func flattenRuleEngineRuleActionsOriginRange(item map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"status": utils.PathSearch("status", item, nil),
		},
	}
}

func flattenRuleEngineRuleActionsAttribute(rawArray, actionsOriginValue []interface{}) []interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for i, item := range rawArray {
		result = append(result, map[string]interface{}{
			"flexible_origin": flattenRuleEngineRuleActionsFlexibleOrigin(utils.PathSearch("flexible_origin",
				item, make([]interface{}, 0)).([]interface{}), utils.PathSearch(fmt.Sprintf("[%d].flexible_origin", i),
				actionsOriginValue, make([]interface{}, 0)).([]interface{})),
			"origin_request_header": flattenRuleEngineRuleActionsOriginRequestHeader(utils.PathSearch("origin_request_header",
				item, make([]interface{}, 0)).([]interface{})),
			"http_response_header": flattenRuleEngineRuleActionsHttpResponseHeader(utils.PathSearch("http_response_header",
				item, make([]interface{}, 0)).([]interface{})),
			"access_control": flattenRuleEngineRuleActionsAccessControl(utils.PathSearch("access_control",
				item, make(map[string]interface{})).(map[string]interface{})),
			"request_limit_rule": flattenRuleEngineRuleActionsRequestLimitRule(utils.PathSearch("request_limit_rule",
				item, make(map[string]interface{})).(map[string]interface{})),
			"origin_request_url_rewrite": flattenRuleEngineRuleActionsOriginRequestUrlRewrite(utils.PathSearch("origin_request_url_rewrite",
				item, make(map[string]interface{})).(map[string]interface{})),
			"cache_rule": flattenRuleEngineRuleActionsCacheRule(utils.PathSearch("cache_rule",
				item, make(map[string]interface{})).(map[string]interface{})),
			"request_url_rewrite": flattenRuleEngineRuleActionsRequestUrlRewrite(utils.PathSearch("request_url_rewrite",
				item, make(map[string]interface{})).(map[string]interface{})),
			"browser_cache_rule": flattenRuleEngineRuleActionsBrowserCacheRule(utils.PathSearch("browser_cache_rule",
				item, make(map[string]interface{})).(map[string]interface{})),
			"error_code_cache": flattenRuleEngineRuleActionsErrorCodeCache(utils.PathSearch("error_code_cache",
				item, make([]interface{}, 0)).([]interface{})),
			"origin_range": flattenRuleEngineRuleActionsOriginRange(utils.PathSearch("origin_range",
				item, make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return result
}

func GetRuleEngineRuleById(client *golangsdk.ServiceClient, domainName string, ruleId string) (interface{}, error) {
	rules, err := listRuleEngineRules(client, domainName)
	if err != nil {
		return nil, err
	}

	rule := utils.PathSearch(fmt.Sprintf("[?rule_id =='%s']|[0]", ruleId), rules, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/cdn/configuration/domains/{domain_name}/rules",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the rule with ID '%s' has been removed", ruleId)),
			},
		}
	}
	return rule, nil
}

func resourceRuleEngineRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		domainName = d.Get("domain_name").(string)
		ruleId     = d.Id()
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	rule, err := GetRuleEngineRuleById(client, domainName, ruleId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CDN rule engine rules")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("domain_name", domainName),
		d.Set("name", utils.PathSearch("name", rule, nil)),
		d.Set("status", utils.PathSearch("status", rule, nil)),
		d.Set("priority", utils.PathSearch("priority", rule, nil)),
		d.Set("conditions", utils.JsonToString(utils.PathSearch("conditions", rule, nil))),
		d.Set("actions", flattenRuleEngineRuleActionsAttribute(utils.PathSearch("actions", rule, make([]interface{}, 0)).([]interface{}),
			d.Get("actions").([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateFullRuleEngineRules(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1.0/cdn/configuration/domains/{domain_name}/rules/{rule_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{domain_name}", d.Get("domain_name").(string))
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRuleEngineRuleBodyParams(d)),
		OkCodes:          []int{204},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceRuleEngineRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	err = updateFullRuleEngineRules(client, d)
	if err != nil {
		return diag.Errorf("error updating CDN rule engine rule: %s", err)
	}

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, []string{"conditions"})
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceRuleEngineRuleRead(ctx, d, meta)
}

func deleteRuleEngineRule(client *golangsdk.ServiceClient, domainName, ruleId string) error {
	httpUrl := "v1.0/cdn/configuration/domains/{domain_name}/rules/{rule_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{domain_name}", domainName)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", ruleId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceRuleEngineRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		domainName = d.Get("domain_name").(string)
		ruleId     = d.Id()
		// The error code of the rule engine rule not found is 'CDN.0001'.
		ruleEngineRuleNotFoundCodes = []string{"CDN.0001"}
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	err = deleteRuleEngineRule(client, domainName, ruleId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.error_code", ruleEngineRuleNotFoundCodes...),
			fmt.Sprintf("error deleting CDN rule engine rule (%s)", ruleId))
	}

	return nil
}

func resourceRuleEngineRuleImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.SplitN(importId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<domain_name>/<id>' or '<domain_name>/<name>', but got '%s'", importId)
	}

	// If the ID is a UUID, set the ID and domain name
	if utils.IsUUID(parts[1]) {
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, d.Set("domain_name", parts[0])
	}
	// If the ID is not a UUID, get the rule by name and set the ID and domain name
	var (
		cfg        = meta.(*config.Config)
		domainName = parts[0]
		ruleName   = parts[1]
	)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}
	rule, err := GetRuleEngineRuleByName(client, domainName, ruleName)
	if err != nil {
		return nil, err
	}
	d.SetId(utils.PathSearch("rule_id", rule, "").(string))

	return []*schema.ResourceData{d}, d.Set("domain_name", parts[0])
}
