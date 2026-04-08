package dli

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

// @API DLI GET /v1/{project_id}/sql-defend-rules
func DataSourceSQLDefendRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSQLDefendRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the SQL defend rules are located.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the queue name used to filter SQL defend rules.`,
			},
			"rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the rule name used to filter SQL defend rules.`,
			},

			// Attributes
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of SQL defend rules that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the rule.`,
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The UUID of the rule.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the rule.`,
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The category of the rule.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the rule.`,
						},
						"engine_rules": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: `The engine rules of the rule.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The project ID of the rule.`,
						},
						"queue_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of queue names that matched filter parameters.`,
						},
						"sys_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The system description of the rule.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the rule.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the rule.`,
						},
					},
				},
			},
		},
	}
}

func buildSQLDefendRulesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("queue_name"); ok {
		res = fmt.Sprintf("%s&queue_name=%v", res, v)
	}
	if v, ok := d.GetOk("rule_name"); ok {
		res = fmt.Sprintf("%s&rule_name=%v", res, v)
	}

	return res
}

func listSQLDefendRules(client *golangsdk.ServiceClient, queryParams string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/sql-defend-rules?limit={limit}"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if queryParams != "" {
		listPath += queryParams
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
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

		rules := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, rules...)
		if len(rules) < limit {
			break
		}

		offset += len(rules)
	}

	return result, nil
}

func flattenSQLDefendRules(rules []interface{}) []map[string]interface{} {
	if len(rules) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"name":         utils.PathSearch("rule_name", rule, nil),
			"uuid":         utils.PathSearch("rule_uuid", rule, nil),
			"id":           utils.PathSearch("rule_id", rule, nil),
			"category":     utils.PathSearch("category", rule, nil),
			"description":  utils.PathSearch("desc", rule, nil),
			"engine_rules": utils.PathSearch("engine_rules", rule, nil),
			"project_id":   utils.PathSearch("project_id", rule, nil),
			"queue_names":  utils.PathSearch("queue_names", rule, make([]interface{}, 0)),
			"sys_desc":     utils.PathSearch("sys_desc", rule, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", rule, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("update_time", rule, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceSQLDefendRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	rules, err := listSQLDefendRules(client, buildSQLDefendRulesQueryParams(d))
	if err != nil {
		return diag.Errorf("error listing SQL defend rules: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenSQLDefendRules(rules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
