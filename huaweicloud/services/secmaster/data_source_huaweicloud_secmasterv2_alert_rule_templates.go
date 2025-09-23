package secmaster

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

// @API SecMaster GET /v2/{project_id}/workspaces/{workspace_id}/siem/alert-rules/templates
func DataSourceAlertRuleTemplatesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlertRuleTemplatesV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accumulated_times": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"alert_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_remediation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_type": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"create_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"custom_properties": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_grouping": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"job_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delay_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"frequency_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"frequency_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"overtime_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"period_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"period_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"simulation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// The field `suppression` is misspelled in the API documentation.
						// The actual names defined in the current schema are valid.
						"suppression": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"triggers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accumulated_times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"expression": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"job_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"update_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time_by_user": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAlertRuleTemplatesv2QueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=1000"

	if v, ok := d.GetOk("template_name"); ok {
		queryParams = fmt.Sprintf("%s&template_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceAlertRuleTemplatesV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v2/{project_id}/workspaces/{workspace_id}/siem/alert-rules/templates"
		result      = make([]interface{}, 0)
		offset      = 0
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath += buildAlertRuleTemplatesv2QueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving alert rule templates (v2): %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		records := utils.PathSearch("records", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		result = append(result, records...)
		offset += len(records)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenAlertRuleTemplatesV2(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAlertRuleTemplatesV2(instancesResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(instancesResp))
	for _, v := range instancesResp {
		rst = append(rst, map[string]interface{}{
			"template_id":         utils.PathSearch("template_id", v, nil),
			"template_name":       utils.PathSearch("template_name", v, nil),
			"accumulated_times":   utils.PathSearch("accumulated_times", v, nil),
			"alert_description":   utils.PathSearch("alert_description", v, nil),
			"alert_name":          utils.PathSearch("alert_name", v, nil),
			"alert_remediation":   utils.PathSearch("alert_remediation", v, nil),
			"alert_type":          utils.PathSearch("alert_type", v, nil),
			"create_by":           utils.PathSearch("create_by", v, nil),
			"create_time":         utils.PathSearch("create_time", v, nil),
			"custom_properties":   utils.PathSearch("custom_properties", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"event_grouping":      utils.PathSearch("event_grouping", v, nil),
			"job_mode":            utils.PathSearch("job_mode", v, nil),
			"process_status":      utils.PathSearch("process_status", v, nil),
			"query":               utils.PathSearch("query", v, nil),
			"query_type":          utils.PathSearch("query_type", v, nil),
			"schedule":            flattenAlertRuleTemplatesV2Schedule(utils.PathSearch("schedule", v, nil)),
			"severity":            utils.PathSearch("severity", v, nil),
			"simulation":          utils.PathSearch("simulation", v, nil),
			"status":              utils.PathSearch("status", v, nil),
			"suppression":         utils.PathSearch("suppression", v, nil),
			"table_name":          utils.PathSearch("table_name", v, nil),
			"triggers":            flattenAlertRuleTemplatesV2Triggers(utils.PathSearch("triggers", v, make([]interface{}, 0)).([]interface{})),
			"update_by":           utils.PathSearch("update_by", v, nil),
			"update_time":         utils.PathSearch("update_time", v, nil),
			"update_time_by_user": utils.PathSearch("update_time_by_user", v, nil),
		})
	}

	return rst
}

func flattenAlertRuleTemplatesV2Schedule(rawSchedule interface{}) []interface{} {
	if rawSchedule == nil {
		return nil
	}

	result := map[string]interface{}{
		"delay_interval":     utils.PathSearch("delay_interval", rawSchedule, nil),
		"frequency_interval": utils.PathSearch("frequency_interval", rawSchedule, nil),
		"frequency_unit":     utils.PathSearch("frequency_unit", rawSchedule, nil),
		"overtime_interval":  utils.PathSearch("overtime_interval", rawSchedule, nil),
		"period_interval":    utils.PathSearch("period_interval", rawSchedule, nil),
		"period_unit":        utils.PathSearch("period_unit", rawSchedule, nil),
	}

	return []interface{}{result}
}

func flattenAlertRuleTemplatesV2Triggers(rawTriggers []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(rawTriggers))
	for _, v := range rawTriggers {
		rst = append(rst, map[string]interface{}{
			"accumulated_times": utils.PathSearch("accumulated_times", v, nil),
			"expression":        utils.PathSearch("expression", v, nil),
			"job_id":            utils.PathSearch("job_id", v, nil),
			"mode":              utils.PathSearch("mode", v, nil),
			"operator":          utils.PathSearch("operator", v, nil),
			"severity":          utils.PathSearch("severity", v, nil),
		})
	}

	return rst
}
