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

// @API Workspace GET /v1/{project_id}/app-center/app-restricted-rules
func DataSourceApplicationRestrictedRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationRestrictedRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the restricted application rules are located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the restricted application rule to be queried.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        restrictedRuleItemSchema(),
				Description: `The list of restricted application rules that match the filter parameters.`,
			},
		},
	}
}

func restrictedRuleItemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the restricted application rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the restricted application rule.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the restricted application rule.`,
			},
			"rule_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source of the restricted application rule.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time of the restricted application rule, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the restricted application rule, in RFC3339 format.`,
			},
			"rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        restrictedRuleDetailSchema(),
				Description: `The detail of the restricted application rule.`,
			},
		},
	}
}

func restrictedRuleDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scope of the restricted application rule.`,
			},
			"product_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        restrictedProductRuleSchema(),
				Description: `The detail of the product rule.`,
			},
			"path_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        restrictedPathRuleSchema(),
				Description: `The detail of the path rule.`,
			},
		},
	}
}

func restrictedProductRuleSchema() *schema.Resource {
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
				Description: `The supported operating system type.`,
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

func restrictedPathRuleSchema() *schema.Resource {
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

func buildRestrictedRulesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	return res
}

func listRestrictedApplicationRules(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-center/app-restricted-rules?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildRestrictedRulesQueryParams(d)

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

		restrictedRules := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, restrictedRules...)
		if len(restrictedRules) < limit {
			break
		}
		offset += len(restrictedRules)
	}

	return result, nil
}

func flattenRestrictedProductRule(productRule interface{}) []interface{} {
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

func flattenRestrictedPathRule(pathRule interface{}) []interface{} {
	if pathRule == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"path": utils.PathSearch("path", pathRule, nil),
		},
	}
}

func flattenRestrictedRuleDetail(ruleDetail interface{}) []interface{} {
	if ruleDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"scope":        utils.PathSearch("scope", ruleDetail, nil),
			"product_rule": flattenRestrictedProductRule(utils.PathSearch("product_rule", ruleDetail, nil)),
			"path_rule":    flattenRestrictedPathRule(utils.PathSearch("path_rule", ruleDetail, nil)),
		},
	}
}

func flattenRestrictedApplicationRules(restrictedRules []interface{}) []interface{} {
	if len(restrictedRules) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(restrictedRules))
	for _, item := range restrictedRules {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", item, nil),
			"name":        utils.PathSearch("name", item, nil),
			"description": utils.PathSearch("description", item, nil),
			"rule_source": utils.PathSearch("rule_source", item, nil),
			"create_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("create_time", item, "").(string))/1000, false),
			"update_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("update_time", item, "").(string))/1000, false),
			"rule": flattenRestrictedRuleDetail(utils.PathSearch("rule", item, nil)),
		})
	}
	return result
}

func dataSourceApplicationRestrictedRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	restrictedRules, err := listRestrictedApplicationRules(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace restricted application rules: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("rules", flattenRestrictedApplicationRules(restrictedRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
