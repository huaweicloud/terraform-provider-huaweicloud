package secmaster

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

// @API SecMaster GET /v2/{project_id}/workspaces/{workspace_id}/siem/alert-rules/templates/{template_id}
func DataSourceAlertRuleTemplateDetailV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlertRuleTemplateDetailV2Read,

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
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_template_id": {
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
			"cu_quota_amount": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_mode_setting": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_overtime_strategy_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"batch_overtime_strategy_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_delay_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"search_delay_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_frequency_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"search_frequency_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_overtime_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"search_overtime_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_period_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"search_period_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_table_name": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"streaming_checkpoint_ttl_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"streaming_checkpoint_ttl_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"streaming_startup_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"streaming_state_ttl_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"job_output_setting": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_custom_properties": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"alert_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_grouping": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"alert_mapping": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"alert_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_remediation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_suppression": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"alert_type": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"entity_extraction": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"field_mapping": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"process_error": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"script": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
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
			"create_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAlertRuleTemplateDetailV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		templateId  = d.Get("template_id").(string)
		httpUrl     = "v2/{project_id}/workspaces/{workspace_id}/siem/alert-rules/templates/{template_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{template_id}", templateId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving alert rule template (v2) detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rule_template_id", utils.PathSearch("template_id", getRespBody, nil)),
		d.Set("template_name", utils.PathSearch("template_name", getRespBody, nil)),
		d.Set("accumulated_times", utils.PathSearch("accumulated_times", getRespBody, nil)),
		d.Set("cu_quota_amount", utils.PathSearch("cu_quota_amount", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("environment", utils.PathSearch("environment", getRespBody, nil)),
		d.Set("job_mode", utils.PathSearch("job_mode", getRespBody, nil)),
		d.Set("job_mode_setting", flattenTemplateV2JobModeSetting(utils.PathSearch("job_mode_setting", getRespBody, nil))),
		d.Set("job_output_setting", flattenTemplateV2JobOutputSetting(utils.PathSearch("job_output_setting", getRespBody, nil))),
		d.Set("process_error", utils.PathSearch("process_error", getRespBody, nil)),
		d.Set("process_status", utils.PathSearch("process_status", getRespBody, nil)),
		d.Set("query_type", utils.PathSearch("query_type", getRespBody, nil)),
		d.Set("script", utils.PathSearch("script", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("table_name", utils.PathSearch("table_name", getRespBody, nil)),
		d.Set("triggers", flattenTemplateV2Triggers(utils.PathSearch("triggers", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("create_by", utils.PathSearch("create_by", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_by", utils.PathSearch("update_by", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTemplateV2JobModeSetting(jobModeSetting interface{}) []interface{} {
	if jobModeSetting == nil {
		return nil
	}

	result := map[string]interface{}{
		"batch_overtime_strategy_interval":  utils.PathSearch("batch_overtime_strategy_interval", jobModeSetting, nil),
		"batch_overtime_strategy_unit":      utils.PathSearch("batch_overtime_strategy_unit", jobModeSetting, nil),
		"search_delay_interval":             utils.PathSearch("search_delay_interval", jobModeSetting, nil),
		"search_delay_unit":                 utils.PathSearch("search_delay_unit", jobModeSetting, nil),
		"search_frequency_interval":         utils.PathSearch("search_frequency_interval", jobModeSetting, nil),
		"search_frequency_unit":             utils.PathSearch("search_frequency_unit", jobModeSetting, nil),
		"search_overtime_interval":          utils.PathSearch("search_overtime_interval", jobModeSetting, nil),
		"search_overtime_unit":              utils.PathSearch("search_overtime_unit", jobModeSetting, nil),
		"search_period_interval":            utils.PathSearch("search_period_interval", jobModeSetting, nil),
		"search_period_unit":                utils.PathSearch("search_period_unit", jobModeSetting, nil),
		"search_table_id":                   utils.PathSearch("search_table_id", jobModeSetting, nil),
		"search_table_name":                 utils.PathSearch("search_table_name", jobModeSetting, nil),
		"streaming_checkpoint_ttl_interval": utils.PathSearch("streaming_checkpoint_ttl_interval", jobModeSetting, nil),
		"streaming_checkpoint_ttl_unit":     utils.PathSearch("streaming_checkpoint_ttl_unit", jobModeSetting, nil),
		"streaming_startup_mode":            utils.PathSearch("streaming_startup_mode", jobModeSetting, nil),
		"streaming_state_ttl_unit":          utils.PathSearch("streaming_state_ttl_unit", jobModeSetting, nil),
	}

	return []interface{}{result}
}

func flattenTemplateV2JobOutputSetting(jobOutputSetting interface{}) []interface{} {
	if jobOutputSetting == nil {
		return nil
	}

	result := map[string]interface{}{
		"alert_custom_properties": utils.PathSearch("alert_custom_properties", jobOutputSetting, nil),
		"alert_description":       utils.PathSearch("alert_description", jobOutputSetting, nil),
		"alert_grouping":          utils.PathSearch("alert_grouping", jobOutputSetting, nil),
		"alert_mapping":           utils.PathSearch("alert_mapping", jobOutputSetting, nil),
		"alert_name":              utils.PathSearch("alert_name", jobOutputSetting, nil),
		"alert_remediation":       utils.PathSearch("alert_remediation", jobOutputSetting, nil),
		"alert_severity":          utils.PathSearch("alert_severity", jobOutputSetting, nil),
		"alert_suppression":       utils.PathSearch("alert_suppression", jobOutputSetting, nil),
		"alert_type":              utils.PathSearch("alert_type", jobOutputSetting, nil),
		"entity_extraction":       utils.PathSearch("entity_extraction", jobOutputSetting, nil),
		"field_mapping":           utils.PathSearch("field_mapping", jobOutputSetting, nil),
	}

	return []interface{}{result}
}

func flattenTemplateV2Triggers(triggers []interface{}) []interface{} {
	if len(triggers) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(triggers))
	for _, v := range triggers {
		result = append(result, map[string]interface{}{
			"accumulated_times": utils.PathSearch("accumulated_times", v, nil),
			"expression":        utils.PathSearch("expression", v, nil),
			"job_id":            utils.PathSearch("job_id", v, nil),
			"mode":              utils.PathSearch("mode", v, nil),
			"operator":          utils.PathSearch("operator", v, nil),
			"severity":          utils.PathSearch("severity", v, nil),
		})
	}

	return result
}
