package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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
				Description: `The region where the app rules are located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the app rule to be queried.`,
			},
			"app_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appRuleSchema(),
				Description: `The list of app rules that match the filter parameters.`,
			},
		},
	}
}

func appRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the app rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the app rule.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the app rule.`,
			},
			"rule_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source of the app rule.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time of the app rule, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the app rule, in RFC3339 format.`,
			},
			"rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        appRuleConfigSchema(),
				Description: `The rule configuration.`,
			},
		},
	}
}

func appRuleConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scope of the rule.`,
			},
			"product_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        productRuleSchema(),
				Description: `The product rule configuration.`,
			},
			"path_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        pathRuleSchema(),
				Description: `The path rule configuration.`,
			},
		},
	}
}

func productRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"identify_condition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The identification condition.`,
			},
			"publisher": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publisher name.`,
			},
			"product_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product name.`,
			},
			"process_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process name.`,
			},
			"support_os": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The supported operating system type.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version number.`,
			},
			"product_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product version number.`,
			},
		},
	}
}

func pathRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The complete path.`,
			},
		},
	}
}

func buildAppRulesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&key=%v", res, v)
	}
	return res
}

func listAppRules(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-center/app-rules?limit={limit}"
		offset  = 0
		limit   = 10
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildAppRulesQueryParams(d)

	opt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%v", listPathWithLimit, strconv.Itoa(offset))
		requestResp, err := client.Request("GET", listPathWithOffset, opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		appRules := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, appRules...)
		if len(appRules) < limit {
			break
		}
		offset += len(appRules)
	}

	return result, nil
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
			"create_time": utils.PathSearch("create_time", item, nil),
			"update_time": utils.PathSearch("update_time", item, nil),
			"rule":        flattenAppRuleConfig(utils.PathSearch("rule", item, nil)),
		})
	}
	return result
}

func flattenAppRuleConfig(rule interface{}) []interface{} {
	if rule == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"scope":        utils.PathSearch("scope", rule, nil),
			"product_rule": flattenProductRule(utils.PathSearch("product_rule", rule, nil)),
			"path_rule":    flattenPathRule(utils.PathSearch("path_rule", rule, nil)),
		},
	}
}

func flattenProductRule(productRule interface{}) []interface{} {
	if productRule == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"identify_condition": utils.PathSearch("identify_condition", productRule, nil),
			"publisher":          utils.PathSearch("publisher", productRule, nil),
			"product_name":       utils.PathSearch("product_name", productRule, nil),
			"process_name":       utils.PathSearch("process_name", productRule, nil),
			"support_os":         utils.PathSearch("support_os", productRule, nil),
			"version":            utils.PathSearch("version", productRule, nil),
			"product_version":    utils.PathSearch("product_version", productRule, nil),
		},
	}
}

func flattenPathRule(pathRule interface{}) []interface{} {
	if pathRule == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"path": utils.PathSearch("path", pathRule, nil),
		},
	}
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
		return diag.Errorf("error querying Workspace app rules: %s", err)
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
