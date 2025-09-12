package cdn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN POST /v1.0/cdn/configuration/domains/{domain_name}/rules
// @API CDN GET /v1.0/cdn/configuration/domains/{domain_name}/rules
// @API CDN PUT /v1.0/cdn/configuration/domains/{domain_name}/rules/full-update
// @API CDN POST /v1.0/cdn/configuration/domains/{domain_name}/rules/batch-update
// @API CDN DELETE /v1.0/cdn/configuration/domains/{domain_name}/rules/{rule_id}

// 不可更新参数列表
var rulesEngineNonUpdatableParams = []string{"domain_name"}

func ResourceRulesEngine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRulesEngineCreate,
		ReadContext:   resourceRulesEngineRead,
		UpdateContext: resourceRulesEngineUpdate,
		DeleteContext: resourceRulesEngineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRulesEngineImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		CustomizeDiff: config.FlexibleForceNew(rulesEngineNonUpdatableParams),
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the CDN rules engine is located.",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The domain name.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule name. The valid length is limit from `1` to `50`.",
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The rule status. Valid values are **on** and **off**.",
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `The rule priority. The valid value is limit from 1 to 100.`,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"conditions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The trigger conditions for the current rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "The rule matching conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"logic": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The logical operator. Valid values are **and** and **or**.",
										ValidateFunc: validation.StringInSlice([]string{"and", "or"}, false),
									},
									"criteria": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The list of matching conditions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_target_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The matching target type.",
													ValidateFunc: validation.StringInSlice([]string{
														"schema", "method", "path", "arg", "extension", "filename",
														"header", "clientip", "clientip_version", "ua", "ngx_variable",
													}, false),
												},
												"match_target_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The matching target name.",
												},
												"match_type": {
													Type:         schema.TypeString,
													Required:     true,
													Description:  "The matching algorithm. Currently only supports **contains**.",
													ValidateFunc: validation.StringInSlice([]string{"contains"}, false),
												},
												"match_pattern": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "The matching content.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"negate": {
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
													Description: "Whether to negate. Defaults to **false**.",
												},
												"case_sensitive": {
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
													Description: "Whether to be case sensitive. Defaults to **false**.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"actions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The actions to be executed after the rule conditions are met.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_control": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The access control configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The access control type. Valid values are **block** and **trust**.",
										ValidateFunc: validation.StringInSlice([]string{"block", "trust"}, false),
									},
								},
							},
						},
						"cache_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The cache rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ttl": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The cache expiration time of resources on CDN nodes.",
									},
									"ttl_unit": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The unit of cache expiration time. Valid values are **s**, **m**, **h**, **d**.",
										ValidateFunc: validation.StringInSlice([]string{"s", "m", "h", "d"}, false),
									},
									"follow_origin": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The source of cache expiration time. Valid values are **on**, **off**, **min_ttl**.",
										ValidateFunc: validation.StringInSlice([]string{"on", "off", "min_ttl"}, false),
									},
									"force_cache": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Whether to enable forced caching. Valid values are **on** and **off**.",
										ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
									},
								},
							},
						},
						"request_url_rewrite": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The request URL rewrite configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redirect_status_code": {
										Type:         schema.TypeInt,
										Optional:     true,
										Description:  "The redirect status code. Valid values are **301**, **302**, **303**, **307**.",
										ValidateFunc: validation.IntInSlice([]int{301, 302, 303, 307}),
									},
									"redirect_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The redirect URL.",
									},
									"redirect_host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The redirect host.",
									},
									"execution_mode": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The execution rule. Valid values are **redirect** and **break**.",
										ValidateFunc: validation.StringInSlice([]string{"redirect", "break"}, false),
									},
								},
							},
						},
						"browser_cache_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The browser cache rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cache_type": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The cache effective type. Valid values are **follow_origin**, **ttl**, **never**.",
										ValidateFunc: validation.StringInSlice([]string{"follow_origin", "ttl", "never"}, false),
									},
									"ttl": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The cache expiration time.",
									},
									"ttl_unit": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "The unit of cache expiration time. Valid values are **s**, **m**, **h**, **d**.",
										ValidateFunc: validation.StringInSlice([]string{"s", "m", "h", "d"}, false),
									},
								},
							},
						},
						"http_response_header": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The HTTP response header configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The HTTP response header parameter name.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The HTTP response header parameter value.",
									},
									"action": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The HTTP response header operation type. Valid values are **set** and **delete**.",
										ValidateFunc: validation.StringInSlice([]string{"set", "delete"}, false),
									},
								},
							},
						},
						"origin_request_header": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The origin request header configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The origin request header parameter name.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The origin request header parameter value.",
									},
									"action": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The origin request header operation type. Valid values are **set** and **delete**.",
										ValidateFunc: validation.StringInSlice([]string{"set", "delete"}, false),
									},
								},
							},
						},
						"request_limit_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The request limit rules configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_rate_after": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The rate limit condition.",
									},
									"limit_rate_value": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The rate limit value.",
									},
								},
							},
						},
						"error_code_cache": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The error code cache configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The error code to be cached.",
									},
									"ttl": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The error code cache time.",
									},
								},
							},
						},
					},
				},
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule ID.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildConditionsMatchBodyParamsV1(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	params := map[string]interface{}{
		"logic": raw["logic"],
	}

	if criteria, ok := raw["criteria"].([]interface{}); ok {
		criteriaParams := make([]map[string]interface{}, 0, len(criteria))
		for _, v := range criteria {
			criteriaRaw := v.(map[string]interface{})
			criteriaParam := map[string]interface{}{
				"match_target_type": criteriaRaw["match_target_type"],
				"match_type":        criteriaRaw["match_type"],
				"match_pattern":     criteriaRaw["match_pattern"],
			}

			if v, ok := criteriaRaw["match_target_name"].(string); ok && v != "" {
				criteriaParam["match_target_name"] = v
			}
			if v, ok := criteriaRaw["negate"].(bool); ok {
				criteriaParam["negate"] = v
			}
			if v, ok := criteriaRaw["case_sensitive"].(bool); ok {
				criteriaParam["case_sensitive"] = v
			}

			criteriaParams = append(criteriaParams, criteriaParam)
		}
		params["criteria"] = criteriaParams
	}

	return params
}

func buildConditionsBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"match": buildConditionsMatchBodyParamsV1(raw["match"].([]interface{})),
	}
}

func buildActionsBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	actionsParams := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		actionRaw := v.(map[string]interface{})
		actionParam := make(map[string]interface{})

		// Access Control
		if accessControl, ok := actionRaw["access_control"].([]interface{}); ok && len(accessControl) > 0 {
			accessControlRaw := accessControl[0].(map[string]interface{})
			actionParam["access_control"] = map[string]interface{}{
				"type": accessControlRaw["type"],
			}
		}

		// Cache Rule
		if cacheRule, ok := actionRaw["cache_rule"].([]interface{}); ok && len(cacheRule) > 0 {
			cacheRuleRaw := cacheRule[0].(map[string]interface{})
			cacheRuleParam := map[string]interface{}{
				"ttl":           cacheRuleRaw["ttl"],
				"ttl_unit":      cacheRuleRaw["ttl_unit"],
				"follow_origin": cacheRuleRaw["follow_origin"],
			}
			if v, ok := cacheRuleRaw["force_cache"].(string); ok && v != "" {
				cacheRuleParam["force_cache"] = v
			}
			actionParam["cache_rule"] = cacheRuleParam
		}

		// Request URL Rewrite
		if requestUrlRewrite, ok := actionRaw["request_url_rewrite"].([]interface{}); ok && len(requestUrlRewrite) > 0 {
			requestUrlRewriteRaw := requestUrlRewrite[0].(map[string]interface{})
			requestUrlRewriteParam := map[string]interface{}{
				"redirect_url":   requestUrlRewriteRaw["redirect_url"],
				"execution_mode": requestUrlRewriteRaw["execution_mode"],
			}
			if v, ok := requestUrlRewriteRaw["redirect_status_code"].(int); ok && v > 0 {
				requestUrlRewriteParam["redirect_status_code"] = v
			}
			if v, ok := requestUrlRewriteRaw["redirect_host"].(string); ok && v != "" {
				requestUrlRewriteParam["redirect_host"] = v
			}
			actionParam["request_url_rewrite"] = requestUrlRewriteParam
		}

		// Browser Cache Rule
		if browserCacheRule, ok := actionRaw["browser_cache_rule"].([]interface{}); ok && len(browserCacheRule) > 0 {
			browserCacheRuleRaw := browserCacheRule[0].(map[string]interface{})
			browserCacheRuleParam := map[string]interface{}{
				"cache_type": browserCacheRuleRaw["cache_type"],
			}
			if v, ok := browserCacheRuleRaw["ttl"].(int); ok && v > 0 {
				browserCacheRuleParam["ttl"] = v
			}
			if v, ok := browserCacheRuleRaw["ttl_unit"].(string); ok && v != "" {
				browserCacheRuleParam["ttl_unit"] = v
			}
			actionParam["browser_cache_rule"] = browserCacheRuleParam
		}

		// HTTP Response Header
		if httpResponseHeader, ok := actionRaw["http_response_header"].([]interface{}); ok && len(httpResponseHeader) > 0 {
			httpResponseHeaderParams := make([]map[string]interface{}, 0, len(httpResponseHeader))
			for _, header := range httpResponseHeader {
				headerRaw := header.(map[string]interface{})
				headerParam := map[string]interface{}{
					"name":   headerRaw["name"],
					"action": headerRaw["action"],
				}
				if v, ok := headerRaw["value"].(string); ok && v != "" {
					headerParam["value"] = v
				}
				httpResponseHeaderParams = append(httpResponseHeaderParams, headerParam)
			}
			actionParam["http_response_header"] = httpResponseHeaderParams
		}

		// Origin Request Header
		if originRequestHeader, ok := actionRaw["origin_request_header"].([]interface{}); ok && len(originRequestHeader) > 0 {
			originRequestHeaderParams := make([]map[string]interface{}, 0, len(originRequestHeader))
			for _, header := range originRequestHeader {
				headerRaw := header.(map[string]interface{})
				headerParam := map[string]interface{}{
					"name":   headerRaw["name"],
					"action": headerRaw["action"],
				}
				if v, ok := headerRaw["value"].(string); ok && v != "" {
					headerParam["value"] = v
				}
				originRequestHeaderParams = append(originRequestHeaderParams, headerParam)
			}
			actionParam["origin_request_header"] = originRequestHeaderParams
		}

		// Request Limit Rules
		if requestLimitRules, ok := actionRaw["request_limit_rules"].([]interface{}); ok && len(requestLimitRules) > 0 {
			requestLimitRulesRaw := requestLimitRules[0].(map[string]interface{})
			actionParam["request_limit_rules"] = map[string]interface{}{
				"limit_rate_after": requestLimitRulesRaw["limit_rate_after"],
				"limit_rate_value": requestLimitRulesRaw["limit_rate_value"],
			}
		}

		// Error Code Cache
		if errorCodeCache, ok := actionRaw["error_code_cache"].([]interface{}); ok && len(errorCodeCache) > 0 {
			errorCodeCacheParams := make([]map[string]interface{}, 0, len(errorCodeCache))
			for _, item := range errorCodeCache {
				itemRaw := item.(map[string]interface{})
				errorCodeCacheParam := map[string]interface{}{
					"code": itemRaw["code"],
					"ttl":  itemRaw["ttl"],
				}
				errorCodeCacheParams = append(errorCodeCacheParams, errorCodeCacheParam)
			}
			actionParam["error_code_cache"] = errorCodeCacheParams
		}

		actionsParams = append(actionsParams, actionParam)
	}

	return actionsParams
}

func buildCreateRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":       d.Get("name"),
		"status":     d.Get("status"),
		"priority":   d.Get("priority"),
		"conditions": buildConditionsBodyParams(d.Get("conditions").([]interface{})),
		"actions":    buildActionsBodyParams(d.Get("actions").([]interface{})),
	}
}

func resourceRulesEngineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	httpUrl := "v1.0/cdn/configuration/domains/{domain_name}/rules"
	client, err := cfg.NewServiceClient("cdnv1", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_name}", domainName)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateRuleBodyParams(d),
		OkCodes:          []int{200, 201, 204},
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CDN rule: %s", err)
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleID, err := getRuleIDByName(client, domainName, d.Get("name").(string))
	if err != nil {
		return diag.Errorf("error getting rule ID after creation: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", domainName, ruleID))

	return resourceRulesEngineRead(ctx, d, meta)
}

func resourceRulesEngineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cdnv1", region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	ruleName := d.Get("name").(string)

	rule, err := GetCdnRule(client, domainName, ruleName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CDN rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rule_id", utils.PathSearch("rule_id", rule, nil)),
		d.Set("name", utils.PathSearch("name", rule, nil)),
		d.Set("status", utils.PathSearch("status", rule, nil)),
		d.Set("priority", utils.PathSearch("priority", rule, nil)),
		d.Set("conditions", flattenCdnRulesEngineConditions(utils.PathSearch("conditions", rule, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("actions", flattenCdnRulesEngineActions(utils.PathSearch("actions", rule, make([]interface{}, 0)).([]interface{}))),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// GetCdnRule is a method used to query CDN rule detail.
func GetCdnRule(client *golangsdk.ServiceClient, domainName, ruleName string) (interface{}, error) {
	httpUrl := "v1.0/cdn/configuration/domains/{domain_name}/rules"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_name}", domainName)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	rules := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
	if len(rules) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	rule := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", ruleName), rules, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func flattenCdnRulesEngineConditions(conditions map[string]interface{}) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"match": flattenCdnRulesEngineMatch(utils.PathSearch("match", conditions, make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenCdnRulesEngineMatch(match map[string]interface{}) []map[string]interface{} {
	if len(match) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"logic":    utils.PathSearch("logic", match, nil),
			"criteria": flattenCdnRulesEngineCriteria(utils.PathSearch("criteria", match, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenCdnRulesEngineCriteria(criteria []interface{}) []map[string]interface{} {
	if len(criteria) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(criteria))
	for _, item := range criteria {
		criteriaItem := item.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"match_target_type": utils.PathSearch("match_target_type", criteriaItem, nil),
			"match_target_name": utils.PathSearch("match_target_name", criteriaItem, nil),
			"match_type":        utils.PathSearch("match_type", criteriaItem, nil),
			"match_pattern":     utils.PathSearch("match_pattern", criteriaItem, make([]interface{}, 0)),
			"negate":            utils.PathSearch("negate", criteriaItem, false),
			"case_sensitive":    utils.PathSearch("case_sensitive", criteriaItem, false),
		})
	}
	return result
}

func flattenCdnRulesEngineActions(actions []interface{}) []map[string]interface{} {
	if len(actions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(actions))
	for _, action := range actions {
		actionItem := action.(map[string]interface{})
		actionMap := make(map[string]interface{})

		// Access Control
		if accessControl := utils.PathSearch("access_control", actionItem, nil); accessControl != nil {
			actionMap["access_control"] = []map[string]interface{}{
				{
					"type": utils.PathSearch("type", accessControl, nil),
				},
			}
		}

		// Cache Rule
		if cacheRule := utils.PathSearch("cache_rule", actionItem, nil); cacheRule != nil {
			actionMap["cache_rule"] = []map[string]interface{}{
				{
					"ttl":           utils.PathSearch("ttl", cacheRule, nil),
					"ttl_unit":      utils.PathSearch("ttl_unit", cacheRule, nil),
					"follow_origin": utils.PathSearch("follow_origin", cacheRule, nil),
					"force_cache":   utils.PathSearch("force_cache", cacheRule, nil),
				},
			}
		}

		// Request URL Rewrite
		if requestUrlRewrite := utils.PathSearch("request_url_rewrite", actionItem, nil); requestUrlRewrite != nil {
			actionMap["request_url_rewrite"] = []map[string]interface{}{
				{
					"redirect_status_code": utils.PathSearch("redirect_status_code", requestUrlRewrite, nil),
					"redirect_url":         utils.PathSearch("redirect_url", requestUrlRewrite, nil),
					"redirect_host":        utils.PathSearch("redirect_host", requestUrlRewrite, nil),
					"execution_mode":       utils.PathSearch("execution_mode", requestUrlRewrite, nil),
				},
			}
		}

		// Browser Cache Rule
		if browserCacheRule := utils.PathSearch("browser_cache_rule", actionItem, nil); browserCacheRule != nil {
			actionMap["browser_cache_rule"] = []map[string]interface{}{
				{
					"cache_type": utils.PathSearch("cache_type", browserCacheRule, nil),
					"ttl":        utils.PathSearch("ttl", browserCacheRule, nil),
					"ttl_unit":   utils.PathSearch("ttl_unit", browserCacheRule, nil),
				},
			}
		}

		// HTTP Response Header
		if httpResponseHeader := utils.PathSearch("http_response_header", actionItem, make([]interface{}, 0)); len(httpResponseHeader.([]interface{})) > 0 {
			headers := make([]map[string]interface{}, 0, len(httpResponseHeader.([]interface{})))
			for _, header := range httpResponseHeader.([]interface{}) {
				headerItem := header.(map[string]interface{})
				headers = append(headers, map[string]interface{}{
					"name":   utils.PathSearch("name", headerItem, nil),
					"value":  utils.PathSearch("value", headerItem, nil),
					"action": utils.PathSearch("action", headerItem, nil),
				})
			}
			actionMap["http_response_header"] = headers
		}

		// Origin Request Header
		if originRequestHeader := utils.PathSearch("origin_request_header", actionItem, make([]interface{}, 0)); len(originRequestHeader.([]interface{})) > 0 {
			headers := make([]map[string]interface{}, 0, len(originRequestHeader.([]interface{})))
			for _, header := range originRequestHeader.([]interface{}) {
				headerItem := header.(map[string]interface{})
				headers = append(headers, map[string]interface{}{
					"name":   utils.PathSearch("name", headerItem, nil),
					"value":  utils.PathSearch("value", headerItem, nil),
					"action": utils.PathSearch("action", headerItem, nil),
				})
			}
			actionMap["origin_request_header"] = headers
		}

		// Request Limit Rules
		if requestLimitRules := utils.PathSearch("request_limit_rules", actionItem, nil); requestLimitRules != nil {
			actionMap["request_limit_rules"] = []map[string]interface{}{
				{
					"limit_rate_after": utils.PathSearch("limit_rate_after", requestLimitRules, nil),
					"limit_rate_value": utils.PathSearch("limit_rate_value", requestLimitRules, nil),
				},
			}
		}

		// Error Code Cache
		if errorCodeCache := utils.PathSearch("error_code_cache", actionItem, nil); errorCodeCache != nil {
			actionMap["error_code_cache"] = []map[string]interface{}{
				{
					"code": utils.PathSearch("code", errorCodeCache, nil),
					"ttl":  utils.PathSearch("ttl", errorCodeCache, nil),
				},
			}
		}

		result = append(result, actionMap)
	}
	return result
}

func resourceRulesEngineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdnv1", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	httpUrl := "v1.0/cdn/configuration/domains/{domain_name}/rules/full-update"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{domain_name}", domainName)
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateRuleBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating CDN rule (%s): %s", d.Get("name").(string), err)
	}
	return resourceRulesEngineRead(ctx, d, meta)
}

func resourceRulesEngineDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdnv1", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	ruleID := d.Get("rule_id").(string)
	httpUrl := "v1.0/cdn/configuration/domains/{domain_name}/rules/{rule_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{domain_name}", domainName)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", ruleID)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting CDN rule (%s): %s", d.Get("name").(string), err)
	}
	return nil
}

func resourceRulesEngineImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import. Format must be <domain_name>/<rule_id>")
	}
	d.Set("domain_name", parts[0])
	d.Set("rule_id", parts[1])
	return []*schema.ResourceData{d}, nil
}

func getRuleIDByName(client *golangsdk.ServiceClient, domainName, ruleName string) (string, error) {
	rule, err := GetCdnRule(client, domainName, ruleName)
	if err != nil {
		return "", fmt.Errorf("error getting CDN rule: %s", err)
	}

	ruleID := utils.PathSearch("rule_id", rule, "").(string)
	if ruleID == "" {
		return "", fmt.Errorf("unable to find the rule ID from the API response")
	}

	return ruleID, nil
}
