package cdn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/configuration/domains/{domain_name}/rules
func DataSourceRuleEngineRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRuleEngineRulesRead,

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The accelerated domain name to which the rule engine rules belong.`,
			},

			// Attributes.
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        ruleEngineRulesAttrRulesElemSchema(),
				Description: `The list of the rule engine rules.`,
			},
		},
	}
}

func ruleEngineRulesAttrRulesElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the rule engine rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the rule engine rule.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether the rule is enabled.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The priority of the rule engine rule.`,
			},
			"conditions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The trigger conditions of the current rule, in JSON format.`,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flexible_origin": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrFlexibleOriginElemSchema(),
							Description: `The list of flexible origin configurations.`,
						},
						"origin_request_header": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrOriginRequestHeaderElemSchema(),
							Description: `The list of origin request header configurations.`,
						},
						"http_response_header": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrHttpResponseHeaderElemSchema(),
							Description: `The list of HTTP response header configurations.`,
						},
						"access_control": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrAccessControlElemSchema(),
							Description: `The access control configuration.`,
						},
						"request_limit_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrRequestLimitRuleElemSchema(),
							Description: `The request rate limit configuration.`,
						},
						"origin_request_url_rewrite": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrOriginRequestUrlRewriteElemSchema(),
							Description: `The origin request URL rewrite configuration.`,
						},
						"cache_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrCacheRuleElemSchema(),
							Description: `The cache rule configuration.`,
						},
						"request_url_rewrite": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrRequestUrlRewriteElemSchema(),
							Description: `The access URL rewrite configuration.`,
						},
						"browser_cache_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrBrowserCacheRuleElemSchema(),
							Description: `The browser cache rule configuration.`,
						},
						"error_code_cache": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrErrorCodeCacheElemSchema(),
							Description: `The list of error code cache configurations.`,
						},
						"origin_range": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        ruleEngineRulesAttrOriginRangeElemSchema(),
							Description: `The origin range configuration.`,
						},
					},
				},
				Description: `The list of actions to be performed when the rules are met.`,
			},
		},
	}
}

func ruleEngineRulesAttrFlexibleOriginElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sources_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source type.`,
			},
			"ip_or_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The origin IP or domain name.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The origin priority.`,
			},
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The origin weight.`,
			},
			"obs_bucket_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OBS bucket type.`,
			},
			"bucket_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The third-party object storage access key.`,
			},
			"bucket_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The third-party object storage region.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The third-party object storage name.`,
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The origin host name.`,
			},
			"origin_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The origin protocol.`,
			},
			"http_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The HTTP port number.`,
			},
			"https_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The HTTPS port number.`,
			},
		},
	}
}

func ruleEngineRulesAttrOriginRequestHeaderElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The back-to-origin request header parameter name.`,
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The back-to-origin request header setting type.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The back-to-origin request header parameter value.`,
			},
		},
	}
}

func ruleEngineRulesAttrHttpResponseHeaderElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HTTP response header parameter name.`,
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The operation type of setting HTTP response header.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HTTP response header parameter value.`,
			},
		},
	}
}

func ruleEngineRulesAttrAccessControlElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The access control type.`,
			},
		},
	}
}

func ruleEngineRulesAttrRequestLimitRuleElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"limit_rate_after": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The rate limit condition.`,
			},
			"limit_rate_value": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The rate limit value.`,
			},
		},
	}
}

func ruleEngineRulesAttrOriginRequestUrlRewriteElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rewrite_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The rewrite type.`,
			},
			"target_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The target URL.`,
			},
			"source_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source URL to be rewritten.`,
			},
		},
	}
}

func ruleEngineRulesAttrCacheRuleElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The cache expiration time.`,
			},
			"ttl_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cache expiration time unit.`,
			},
			"follow_origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cache expiration time source.`,
			},
			"force_cache": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether to enable forced caching.`,
			},
		},
	}
}

func ruleEngineRulesAttrRequestUrlRewriteElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"redirect_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The redirect URL.`,
			},
			"execution_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution mode.`,
			},
			"redirect_status_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The redirect status code.`,
			},
			"redirect_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The redirect host.`,
			},
		},
	}
}

func ruleEngineRulesAttrBrowserCacheRuleElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cache_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cache effective type.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The cache expiration time.`,
			},
			"ttl_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cache expiration time unit.`,
			},
		},
	}
}

func ruleEngineRulesAttrErrorCodeCacheElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The error code to be cached.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The error code cache time.`,
			},
		},
	}
}

func ruleEngineRulesAttrOriginRangeElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The origin range status.`,
			},
		},
	}
}

func flattenRuleEngineRules(rules []interface{}) []map[string]interface{} {
	if len(rules) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("rule_id", rule, nil),
			"name":       utils.PathSearch("name", rule, nil),
			"status":     utils.PathSearch("status", rule, nil),
			"priority":   utils.PathSearch("priority", rule, nil),
			"conditions": utils.JsonToString(utils.PathSearch("conditions", rule, nil)),
			"actions": flattenRuleEngineRuleActionsAttribute(utils.PathSearch("actions", rule,
				make([]interface{}, 0)).([]interface{}), nil),
		})
	}

	return result
}

func dataSourceRuleEngineRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		domainName = d.Get("domain_name").(string)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	rules, err := listRuleEngineRules(client, domainName)
	if err != nil {
		return diag.Errorf("error querying CDN rule engine rules: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("rules", flattenRuleEngineRules(rules)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
