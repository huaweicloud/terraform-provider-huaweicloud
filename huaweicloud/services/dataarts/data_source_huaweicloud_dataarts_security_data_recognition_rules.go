package dataarts

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

// @API DataArtsStudio GET /v1/{project_id}/security/data-classification/rule
func DataSourceSecurityDataRecognitionRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityDataRecognitionRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the data recognition rules are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the data recognition rules belong.`,
			},

			// Optional parameters.
			"secrecy_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The secrecy level to which the data recognition rules belong.`,
			},
			"rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the specified data recognition rule to be queried.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the data recognition rules to be queried.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the data recognition rules are enabled to be queried.`,
			},

			// Attributes.
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of data recognition rules that matched filter parameters.`,
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
						"combine_expression": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The combine expression of the data recognition rule.`,
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
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID to which the data recognition rule belongs.`,
						},
						"match_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The match type of the data recognition rule.`,
						},
					},
				},
			},
		},
	}
}

func buildSecurityDataRecognitionRulesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("secrecy_level"); ok {
		res = fmt.Sprintf("%s&secrecy_level=%v", res, v)
	}
	if v, ok := d.GetOk("rule_name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, v)
	}
	if v, ok := d.GetOk("enable"); ok {
		res = fmt.Sprintf("%s&enable=%v", res, v)
	}

	return res
}

func listSecurityDataRecognitionRules(client *golangsdk.ServiceClient, workspaceId string, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/data-classification/rule?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		rules := utils.PathSearch("content", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, rules...)
		if len(rules) < limit {
			break
		}
		offset += len(rules)
	}

	return result, nil
}

func flattenSecurityDataRecognitionRules(rules []interface{}) []map[string]interface{} {
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
			"combine_expression": utils.PathSearch("combine_expression", rule, nil),
			"description":        utils.PathSearch("description", rule, nil),
			"builtin_rule_id":    utils.PathSearch("builtin_rule_id", rule, nil),
			"category_id":        utils.PathSearch("category_id", rule, nil),
			"instance_id":        utils.PathSearch("instance_id", rule, nil),
			"match_type":         utils.PathSearch("match_type", rule, nil),
			"created_by":         utils.PathSearch("created_by", rule, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("created_at", rule, float64(0)).(float64))/1000, false),
			"updated_by": utils.PathSearch("updated_by", rule, nil),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("updated_at", rule, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceSecurityDataRecognitionRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	rules, err := listSecurityDataRecognitionRules(client, workspaceId, buildSecurityDataRecognitionRulesQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying DataArts Security data recognition rules: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenSecurityDataRecognitionRules(rules)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
