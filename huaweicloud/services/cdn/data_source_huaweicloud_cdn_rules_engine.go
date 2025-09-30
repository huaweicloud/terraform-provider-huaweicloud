package cdn

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

// @API CDN GET /v1.0/cdn/configuration/domains/{domain_name}/rules
func DataSourceRulesEngine() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRulesEngineRead,

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
				Description: "The domain name to query rules for.",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of CDN rules that match the filter parameters.",
				Elem:        ruleSchema(),
			},
		},
	}
}

func ruleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule name.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule status, on: enabled, off: disabled.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The rule priority, the higher the value, the higher the priority.",
			},
			"conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The trigger conditions for the current rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The rule matching conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"logic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The logical operator.",
									},
									"criteria": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The list of matching conditions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_target_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The matching target type.",
												},
												"match_target_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The matching target name.",
												},
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The matching algorithm.",
												},
												"match_pattern": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The matching content.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"negate": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to negate.",
												},
												"case_sensitive": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to be case sensitive.",
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
				Computed:    true,
				Description: "The actions to be executed after the rule conditions are met.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_control": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The access control configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The access control type.",
									},
								},
							},
						},
						"cache_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The cache rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ttl": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The cache expiration time of resources on CDN nodes.",
									},
									"ttl_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unit of cache expiration time.",
									},
									"follow_origin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The source of cache expiration time.",
									},
									"force_cache": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to enable forced caching.",
									},
								},
							},
						},
						"request_url_rewrite": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The request URL rewrite configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redirect_status_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The redirect status code.",
									},
									"redirect_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The redirect URL.",
									},
									"redirect_host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The redirect host.",
									},
									"execution_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The execution rule.",
									},
								},
							},
						},
						"browser_cache_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The browser cache rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cache_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cache effective type.",
									},
									"ttl": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The cache expiration time.",
									},
									"ttl_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unit of cache expiration time.",
									},
								},
							},
						},
						"http_response_header": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The HTTP response header configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The HTTP response header parameter name.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The HTTP response header parameter value.",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The HTTP response header operation type.",
									},
								},
							},
						},
						"origin_request_header": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The origin request header configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The origin request header parameter name.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The origin request header parameter value.",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The origin request header operation type.",
									},
								},
							},
						},
						"request_limit_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The request limit rules configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_rate_after": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The rate limit condition.",
									},
									"limit_rate_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The rate limit value.",
									},
								},
							},
						},
						"error_code_cache": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The error code cache configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The error code to be cached.",
									},
									"ttl": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The error code cache time.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func listRules(client *golangsdk.ServiceClient, domainName string) ([]interface{}, error) {
	var (
		httpUrl = "v1.0/cdn/configuration/domains/{domain_name}/rules"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{domain_name}", domainName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	rules := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
	result = append(result, rules...)

	return result, nil
}

func flattenRulesEngine(rules []interface{}) []interface{} {
	if len(rules) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		ruleItem := rule.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"rule_id":    utils.PathSearch("rule_id", ruleItem, nil),
			"name":       utils.PathSearch("name", ruleItem, nil),
			"status":     utils.PathSearch("status", ruleItem, nil),
			"priority":   utils.PathSearch("priority", ruleItem, nil),
			"conditions": flattenRulesEngineConditions(utils.PathSearch("conditions", ruleItem, make(map[string]interface{})).(map[string]interface{})),
			"actions":    flattenRulesEngineActions(utils.PathSearch("actions", ruleItem, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenRulesEngineConditions(conditions map[string]interface{}) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"match": flattenRulesEngineMatch(utils.PathSearch("match", conditions, make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenRulesEngineMatch(match map[string]interface{}) []map[string]interface{} {
	if len(match) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"logic":    utils.PathSearch("logic", match, nil),
			"criteria": flattenRulesEngineCriteria(utils.PathSearch("criteria", match, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenRulesEngineCriteria(criteria []interface{}) []map[string]interface{} {
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

func flattenRulesEngineActions(actions []interface{}) []map[string]interface{} {
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

func dataSourceRulesEngineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cdnv1", region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	rules, err := listRules(client, domainName)
	if err != nil {
		return diag.Errorf("error querying CDN rules: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("rules", flattenRulesEngine(rules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
