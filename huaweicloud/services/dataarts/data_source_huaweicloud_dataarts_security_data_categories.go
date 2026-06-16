package dataarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/data-category
func DataSourceSecurityDataCategories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityDataCategoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the data categories are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the data categories belong.`,
			},

			// Attributes.
			"categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of data categories.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data category.`,
						},
						"category_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the data category.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the data category.`,
						},
						"category_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The level of the data category in the tree.`,
						},
						"root_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the root node of the category tree.`,
						},
						"parent_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the parent category.`,
						},
						"category_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The path of the category in the tree.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID to which the data category belongs.`,
						},
						"synchronize": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the data category is synchronized with assets.`,
						},
						"children": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The children information of the data category, in JSON format.`,
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of data classification rules associated with the category.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the rule.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the rule.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the rule.`,
									},
									"secrecy_level": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The secrecy level name.`,
									},
									"secrecy_level_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The secrecy level number.`,
									},
									"enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the rule is enabled.`,
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The method of the rule.`,
									},
									"content_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The content expression of the rule.`,
									},
									"column_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The column expression of the rule.`,
									},
									"comment_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The comment expression of the rule.`,
									},
									"combine_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The combine expression of the rule.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The description of the rule.`,
									},
									"builtin_rule_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The builtin rule ID.`,
									},
									"match_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The match type of the rule.`,
									},
									"created_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The creator of the rule.`,
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The creation time of the rule, in RFC3339 format.`,
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The updater of the rule.`,
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The latest update time of the rule, in RFC3339 format.`,
									},
								},
							},
						},
						"create_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the data category.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the data category, in RFC3339 format.`,
						},
						"update_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The updater of the data category.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the data category, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func listSecurityDataCategories(client *golangsdk.ServiceClient, workspaceId string) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/security/data-category"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("category_groups", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceSecurityDataCategoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	categories, err := listSecurityDataCategories(client, workspaceId)
	if err != nil {
		return diag.Errorf("error querying DataArts Security data categories: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("categories", flattenSecurityDataCategories(categories)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityDataCategories(categories []interface{}) []interface{} {
	if len(categories) < 1 {
		return []interface{}{}
	}

	result := make([]interface{}, 0, len(categories))
	for _, v := range categories {
		result = append(result, map[string]interface{}{
			"category_id":    utils.PathSearch("category_id", v, nil),
			"category_name":  utils.PathSearch("category_name", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"category_level": int(utils.PathSearch("category_level", v, float64(0)).(float64)),
			"root_id":        utils.PathSearch("root_id", v, nil),
			"parent_id":      utils.PathSearch("parent_id", v, nil),
			"category_path":  utils.PathSearch("category_path", v, nil),
			"instance_id":    utils.PathSearch("instance_id", v, nil),
			"synchronize":    utils.PathSearch("synchronize", v, nil),
			"children":       utils.JsonToString(utils.PathSearch("children", v, nil)),
			"rules": flattenSecurityDataCategoryRules(utils.PathSearch("rule_list",
				v, make([]interface{}, 0)).([]interface{})),
			"create_by": utils.PathSearch("create_by", v, nil),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				v, float64(0)).(float64))/1000, false),
			"update_by": utils.PathSearch("update_by", v, nil),
			"update_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
				v, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func flattenSecurityDataCategoryRules(rules []interface{}) []interface{} {
	if len(rules) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("uuid", rule, nil),
			"name":               utils.PathSearch("name", rule, nil),
			"type":               utils.PathSearch("rule_type", rule, nil),
			"secrecy_level":      utils.PathSearch("secrecy_level", rule, nil),
			"secrecy_level_num":  utils.PathSearch("secrecy_level_num", rule, nil),
			"enable":             utils.PathSearch("enable", rule, nil),
			"method":             utils.PathSearch("method", rule, nil),
			"content_expression": utils.PathSearch("content_expression", rule, nil),
			"column_expression":  utils.PathSearch("column_expression", rule, nil),
			"comment_expression": utils.PathSearch("commit_expression", rule, nil),
			"combine_expression": utils.PathSearch("combine_expression", rule, nil),
			"description":        utils.PathSearch("description", rule, nil),
			"builtin_rule_id":    utils.PathSearch("builtin_rule_id", rule, nil),
			"match_type":         utils.PathSearch("match_type", rule, nil),
			"created_by":         utils.PathSearch("created_by", rule, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("created_at",
				rule, float64(0)).(float64))/1000, false),
			"updated_by": utils.PathSearch("updated_by", rule, nil),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("updated_at",
				rule, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
