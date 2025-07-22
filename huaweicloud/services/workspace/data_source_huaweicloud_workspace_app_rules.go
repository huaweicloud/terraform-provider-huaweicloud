package workspace

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/app-center/app-rules
func DataSourceAppRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the application rules are located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the application rule to be queried.`,
			},
			"app_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appRuleItemSchema(),
				Description: `The list of application rules that match the filter parameters.`,
			},
		},
	}
}

func appRuleItemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the application rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application rule.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the application rule.`,
			},
			"rule_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source of the application rule.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time of the application rule, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the application rule, in RFC3339 format.`,
			},
			"rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appRuleSchema(),
				Description: `The application rule configuration.`,
			},
		},
	}
}

func appRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scope of the application rule.`,
			},
			"product_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appProductRule(),
				Description: `The detail of the product rule.`,
			},
			"path_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appPathRule(),
				Description: `The detail of the path rule.`,
			},
		},
	}
}

func appProductRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"identify_condition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The identify condition of the product rule.`,
			},
			"publisher": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publisher of the product.`,
			},
			"product_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the product.`,
			},
			"process_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process name of the product.`,
			},
			"support_os": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The list of the supported operating system types.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the product rule.`,
			},
			"product_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the product.`,
			},
		},
	}
}

func appPathRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The path where the product is installed.`,
			},
		},
	}
}

func flattenAppRules(appRules []interface{}) []interface{} {
	if len(appRules) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(appRules))
	for _, item := range appRules {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", item, nil),
			"name":        utils.PathSearch("name", item, nil),
			"description": utils.PathSearch("description", item, nil),
			"rule_source": utils.PathSearch("rule_source", item, nil),
			"create_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("create_time", item, "").(string))/1000, false),
			"update_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("update_time", item, "").(string))/1000, false),
			"rule": flattenAppRuleConfig(utils.PathSearch("rule", item, nil)),
		})
	}
	return result
}

func dataSourceAppRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	appRules, err := listAppRules(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace application rules: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("app_rules", flattenAppRules(appRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
