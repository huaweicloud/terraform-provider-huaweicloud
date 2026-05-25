package dataarts

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/data-classification/rule/group
func DataSourceSecurityDataRecognitionRuleGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityDataRecognitionRuleGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the data recognition rule groups are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the data recognition rule groups belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the data recognition rule groups to be queried.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the data recognition rule groups to be queried.`,
			},

			// Attributes.
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of data recognition rule groups that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data recognition rule group.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the data recognition rule group.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the data recognition rule group.`,
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the data recognition rule group.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the data recognition rule group, in RFC3339 format.`,
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The updater of the data recognition rule group.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the data recognition rule group, in RFC3339 format.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The project ID to which the data recognition rule group belongs.`,
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of data recognition rules that the group contains.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the data recognition rule.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the data recognition rule.`,
									},
									"secrecy_level": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The secrecy level name of the data recognition rule.`,
									},
									"secrecy_level_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The secrecy level number of the data recognition rule.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the data recognition rule.`,
									},
									"enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the data recognition rule is enabled.`,
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The method of the data recognition rule.`,
									},
									"content_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The content expression of the data recognition rule.`,
									},
									"column_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The column expression of the data recognition rule.`,
									},
									"comment_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The comment expression of the data recognition rule.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The description of the data recognition rule.`,
									},
									"created_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The creator of the data recognition rule.`,
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The creation time of the data recognition rule, in RFC3339 format.`,
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The updater of the data recognition rule.`,
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The latest update time of the data recognition rule, in RFC3339 format.`,
									},
									"builtin_rule_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The builtin rule ID of the data recognition rule.`,
									},
									"category_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The category ID to which the data recognition rule belongs.`,
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

func buildSecurityDataRecognitionRuleGroupsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, v)
	}

	return res
}

func flattenSecurityDataRecognitionRuleGroupsRules(rules []interface{}) []map[string]interface{} {
	if len(rules) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("uuid", rule, nil),
			"type":               utils.PathSearch("rule_type", rule, nil),
			"secrecy_level":      utils.PathSearch("secrecy_level", rule, nil),
			"secrecy_level_num":  int(utils.PathSearch("secrecy_level_num", rule, float64(0)).(float64)),
			"name":               utils.PathSearch("name", rule, nil),
			"enable":             utils.PathSearch("enable", rule, nil),
			"method":             utils.PathSearch("method", rule, nil),
			"content_expression": utils.PathSearch("content_expression", rule, nil),
			"column_expression":  utils.PathSearch("column_expression", rule, nil),
			"comment_expression": utils.PathSearch("commit_expression", rule, nil),
			"description":        utils.PathSearch("description", rule, nil),
			"created_by":         utils.PathSearch("created_by", rule, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("created_at", rule, float64(0)).(float64))/1000, false),
			"updated_by": utils.PathSearch("updated_by", rule, nil),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("updated_at", rule, float64(0)).(float64))/1000, false),
			"builtin_rule_id": utils.PathSearch("builtin_rule_id", rule, nil),
			"category_id":     utils.PathSearch("category_id", rule, nil),
		})
	}

	return result
}

func flattenSecurityDataRecognitionRuleGroups(groups []interface{}) []map[string]interface{} {
	if len(groups) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		rulesRaw := utils.PathSearch("rules", group, nil)
		rules := make([]interface{}, 0)
		if rulesRaw != nil {
			rules = rulesRaw.([]interface{})
		}

		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("uuid", group, nil),
			"name":        utils.PathSearch("name", group, nil),
			"description": utils.PathSearch("description", group, nil),
			"created_by":  utils.PathSearch("created_by", group, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("created_at", group, float64(0)).(float64))/1000, false),
			"updated_by": utils.PathSearch("updated_by", group, nil),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("updated_at", group, float64(0)).(float64))/1000, false),
			"project_id": utils.PathSearch("project_id", group, nil),
			"rules":      flattenSecurityDataRecognitionRuleGroupsRules(rules),
		})
	}

	return result
}

func dataSourceSecurityDataRecognitionRuleGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	groups, err := listSecurityDataRecognitionRuleGroups(client, workspaceId, buildSecurityDataRecognitionRuleGroupsQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying DataArts Security data recognition rule groups: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenSecurityDataRecognitionRuleGroups(groups)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
