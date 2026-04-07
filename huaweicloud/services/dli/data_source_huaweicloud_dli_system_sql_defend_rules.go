package dli

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

// @API DLI GET /v1/{project_id}/sql-defend-sys-rules
func DataSourceSystemSQLDefendRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSystemSQLDefendRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the system SQL defend rules are located.`,
			},

			// Attributes
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of system SQL defend rules that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the rule type.`,
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The category of the rule.`,
						},
						"engines": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of supported engines that matched filter parameters.`,
						},
						"actions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of executable actions that matched filter parameters.`,
						},
						"no_limit": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the rule has a limit value.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the rule.`,
						},
						"param": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of rule parameters that matched filter parameters.`,
							Elem:        systemSQLDefendRuleParamElem(),
						},
					},
				},
			},
		},
	}
}

func systemSQLDefendRuleParamElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_value": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The default value of the threshold.`,
			},
			"min": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum value of the threshold.`,
			},
			"max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum value of the threshold.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the parameter.`,
			},
		},
	}
}

func flattenSystemSQLDefendRules(rules []interface{}) []map[string]interface{} {
	if len(rules) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("rule_id", rule, nil),
			"category":    utils.PathSearch("category", rule, nil),
			"engines":     utils.PathSearch("engines", rule, make([]interface{}, 0)),
			"actions":     utils.PathSearch("actions", rule, make([]interface{}, 0)),
			"no_limit":    utils.PathSearch("no_limit", rule, false),
			"description": utils.PathSearch("desc", rule, nil),
			"param":       flattenRuleParam(utils.PathSearch("param", rule, nil)),
		})
	}

	return result
}

func flattenRuleParam(paramRaw interface{}) []map[string]interface{} {
	if paramRaw == nil || paramRaw == "" {
		return nil
	}

	return []map[string]interface{}{
		{
			"default_value": utils.PathSearch("defaultValue", paramRaw, nil),
			"min":           utils.PathSearch("min", paramRaw, nil),
			"max":           utils.PathSearch("max", paramRaw, nil),
			"description":   utils.PathSearch("desc", paramRaw, nil),
		},
	}
}

func dataSourceSystemSQLDefendRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	httpUrl := "v1/{project_id}/sql-defend-sys-rules"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving system SQL defend rules: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	rules := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenSystemSQLDefendRules(rules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
