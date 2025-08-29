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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/alert-rules/templates/{template_id}
func DataSourceAlertRuleTemplateDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlertRuleTemplateDetailRead,

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
			"data_source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
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
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"event_grouping": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"frequency_unit": {
							Type:     schema.TypeString,
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
						"delay_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"overtime_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"triggers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accumulated_times": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlertRuleTemplateDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		templateId  = d.Get("template_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/alert-rules/templates/{template_id}"
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
		return diag.Errorf("error retrieving alert rule template detail: %s", err)
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
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
		d.Set("template_name", utils.PathSearch("template_name", getRespBody, nil)),
		d.Set("data_source", utils.PathSearch("data_source", getRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getRespBody, nil)),
		d.Set("query", utils.PathSearch("query", getRespBody, nil)),
		d.Set("query_type", utils.PathSearch("query_type", getRespBody, nil)),
		d.Set("severity", utils.PathSearch("severity", getRespBody, nil)),
		d.Set("custom_properties", utils.PathSearch("custom_properties", getRespBody, nil)),
		d.Set("event_grouping", utils.PathSearch("event_grouping", getRespBody, nil)),
		d.Set("schedule", flattenTemplateSchedule(utils.PathSearch("schedule", getRespBody, nil))),
		d.Set("triggers", flattenTemplateTriggers(utils.PathSearch("triggers", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTemplateSchedule(schedule interface{}) []interface{} {
	if schedule == nil {
		return nil
	}

	result := map[string]interface{}{
		"frequency_interval": utils.PathSearch("frequency_interval", schedule, nil),
		"frequency_unit":     utils.PathSearch("frequency_unit", schedule, nil),
		"period_interval":    utils.PathSearch("period_interval", schedule, nil),
		"period_unit":        utils.PathSearch("period_unit", schedule, nil),
		"delay_interval":     utils.PathSearch("delay_interval", schedule, nil),
		"overtime_interval":  utils.PathSearch("overtime_interval", schedule, nil),
	}

	return []interface{}{result}
}

func flattenTemplateTriggers(triggers []interface{}) []interface{} {
	if len(triggers) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(triggers))
	for _, v := range triggers {
		result = append(result, map[string]interface{}{
			"mode":              utils.PathSearch("mode", v, nil),
			"operator":          utils.PathSearch("operator", v, nil),
			"expression":        utils.PathSearch("expression", v, nil),
			"severity":          utils.PathSearch("severity", v, nil),
			"accumulated_times": utils.PathSearch("accumulated_times", v, nil),
		})
	}

	return result
}
