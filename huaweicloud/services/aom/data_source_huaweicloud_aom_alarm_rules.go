package aom

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

// @API AOM GET /v4/{project_id}/alarm-rules
func DataSourceAlarmRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alarm_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_rule_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_rule_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prom_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bind_notification_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"related_cce_clusters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prom_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_alarm_spec":   dataSourceSchemaMetricAlarmSpec(),
						"event_alarm_spec":    dataSourceSchemeEventAlarmSpec(),
						"alarm_notifications": dataSourceSchemeAlarmNotifications(),
						"status": {
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

func dataSourceSchemaMetricAlarmSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"monitor_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"no_data_conditions":  dataSourceSchemeNoDataConditions(),
				"alarm_tags":          dataSourceSchemeAlarmTags(),
				"trigger_conditions":  dataSourceSchemeMetricTriggerConditions(),
				"monitor_objects":     dataSourceSchemeMonitorObjects(),
				"recovery_conditions": dataSourceSchemeRecoveryConditions(),
			},
		},
	}
}

func dataSourceSchemeEventAlarmSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alarm_source": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"event_source": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"monitor_objects":    dataSourceSchemeMonitorObjects(),
				"trigger_conditions": dataSourceSchemeEventTriggerConditions(),
			},
		},
	}
}

func dataSourceSchemeNoDataConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"no_data_timeframe": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"no_data_alert_state": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"notify_no_data": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceSchemeAlarmTags() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_tags": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"custom_tags": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"custom_annotations": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceSchemeMetricTriggerConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"metric_query_mode": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"promql": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"aggregation_window": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"query_match": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"aggregate_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_labels": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"aggregation_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"operator": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"thresholds": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"trigger_times": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"trigger_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"trigger_interval": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"expression": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"mix_promql": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_statistic_method": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"query_param": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_namespace": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_unit": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"promql_expr": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"promql_for": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"aom_monitor_level": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceSchemeEventTriggerConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"event_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"trigger_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"aggregation_window": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"operator": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"thresholds": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeInt},
				},
				"frequency": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceSchemeRecoveryConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"recovery_timeframe": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceSchemeMonitorObjects() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString}},
	}
}

func dataSourceSchemeAlarmNotifications() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"notification_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"route_group_enable": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"route_group_rule": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"notification_enable": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"bind_notification_rule_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"notify_resolved": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"notify_triggered": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"notify_frequency": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceAlarmRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	results, err := getAlarmRules(cfg, client, d)
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
		d.Set("alarm_rules", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAlarmRules(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v4/{project_id}/alarm-rules"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListAlarmRulesQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeadersForDataSource(cfg, d),
	}

	// if limit and offset is empty, return all
	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the rules list: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening the rules list: %s", err)
	}

	// there are two different kinds of returns
	rules, ok := listRespBody.([]interface{})
	if !ok {
		rules = utils.PathSearch("alarm_rules", listRespBody, make([]interface{}, 0)).([]interface{})
	}

	results := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		results = append(results, flattenAlarmRule(rule))
	}

	return results, nil
}

func flattenAlarmRule(rule interface{}) map[string]interface{} {
	rst := map[string]interface{}{
		"name":                  utils.PathSearch("alarm_rule_name", rule, nil),
		"alarm_rule_id":         strconv.FormatFloat(utils.PathSearch("alarm_rule_id", rule, float64(0)).(float64), 'f', -1, 64),
		"type":                  utils.PathSearch("alarm_rule_type", rule, nil),
		"prom_instance_id":      utils.PathSearch("prom_instance_id", rule, nil),
		"enable":                utils.PathSearch("alarm_rule_enable", rule, nil),
		"status":                utils.PathSearch("alarm_rule_status", rule, nil),
		"description":           utils.PathSearch("alarm_rule_description", rule, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", rule, nil),
		"created_at": utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("alarm_create_time", rule, float64(0)).(float64))/1000, true),
		"updated_at": utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("alarm_update_time", rule, float64(0)).(float64))/1000, true),
		"metric_alarm_spec":   flattenV4MetricAlarmSpec(utils.PathSearch("metric_alarm_spec", rule, nil)),
		"event_alarm_spec":    flattenV4EventAlarmSpec(utils.PathSearch("event_alarm_spec", rule, nil)),
		"alarm_notifications": flattenV4AlarmNotifications(utils.PathSearch("alarm_notifications", rule, nil)),
	}
	return rst
}

func buildListAlarmRulesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("alarm_rule_name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("event_source"); ok {
		res = fmt.Sprintf("%s&event_source=%v", res, v)
	}
	if v, ok := d.GetOk("event_severity"); ok {
		res = fmt.Sprintf("%s&event_severity=%v", res, v)
	}
	if v, ok := d.GetOk("alarm_rule_type"); ok {
		res = fmt.Sprintf("%s&alarm_rule_type=%v", res, v)
	}
	if v, ok := d.GetOk("alarm_rule_status"); ok {
		res = fmt.Sprintf("%s&alarm_rule_status=%v", res, v)
	}
	if v, ok := d.GetOk("prom_instance_id"); ok {
		res = fmt.Sprintf("%s&prom_instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("bind_notification_rule_id"); ok {
		res = fmt.Sprintf("%s&bind_notification_rule_id=%v", res, v)
	}
	if v, ok := d.GetOk("related_cce_clusters"); ok {
		res = fmt.Sprintf("%s&related_cce_clusters=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
