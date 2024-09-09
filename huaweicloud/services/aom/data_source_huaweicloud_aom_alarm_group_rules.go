package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM GET /v2/{project_id}/alert/group-rules
func DataSourceAlarmGroupRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmGroupRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detail": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bind_notification_rule_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"match": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operate": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"group_by": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"group_wait": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_repeat_waiting": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlarmGroupRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	results, err := listAlarmGroupRules(cfg, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenGroupRules(results.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listAlarmGroupRules(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/alert/group-rules"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeadersForDataSource(cfg, d),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving alarm group rules: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening alarm group rules: %s", err)
	}

	return listRespBody, nil
}

func flattenGroupRules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"name":                  utils.PathSearch("name", rule, nil),
			"detail":                flattenAlarmGroupRuleDetail(utils.PathSearch("detail", rule, make([]interface{}, 0)).([]interface{})),
			"group_by":              utils.PathSearch("group_by", rule, nil),
			"group_wait":            utils.PathSearch("group_wait", rule, nil),
			"group_interval":        utils.PathSearch("group_interval", rule, nil),
			"group_repeat_waiting":  utils.PathSearch("group_repeat_waiting", rule, nil),
			"description":           utils.PathSearch("desc", rule, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", rule, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", rule, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("update_time", rule, float64(0)).(float64))/1000, false),
		})
	}
	return result
}
