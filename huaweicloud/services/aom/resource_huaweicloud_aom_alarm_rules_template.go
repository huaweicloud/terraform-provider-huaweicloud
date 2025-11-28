package aom

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v4/{project_id}/alarm-rules-template
// @API AOM PUT /v4/{project_id}/alarm-rules-template
// @API AOM GET /v4/{project_id}/alarm-rules-template
// @API AOM DELETE /v4/{project_id}/alarm-rules-template
func ResourceAlarmRulesTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmRulesTemplateCreate,
		ReadContext:   resourceAlarmRulesTemplateRead,
		UpdateContext: resourceAlarmRulesTemplateUpdate,
		DeleteContext: resourceAlarmRulesTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_template_spec_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"related_cloud_service": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"related_cce_clusters": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"related_prometheus_instances": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"alarm_notification": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     resourceSchemeV4AlarmNotification(),
						},
						"alarm_template_spec_items": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alarm_rule_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"alarm_rule_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"alarm_rule_description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"event_alarm_spec":  resourceSchemeEventAlarmTemplateSpec(),
									"metric_alarm_spec": resourceSchemeMetricAlarmTemplateSpec(),
								},
							},
						},
					},
				},
			},
			"templating": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"list": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"query": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
	}
}

func resourceSchemeEventAlarmTemplateSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alarm_subtype": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"alarm_source": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"event_source": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"trigger_conditions": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"trigger_type": {
								Type:     schema.TypeString,
								Required: true,
							},
							"event_name": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"thresholds": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeInt},
							},
							"aggregation_window": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"operator": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"frequency": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
				"monitor_objects": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeMap,
						Elem: &schema.Schema{Type: schema.TypeString},
					},
				},
				"monitor_object_templates": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func resourceSchemeMetricAlarmTemplateSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alarm_subtype": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"alarm_source": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"monitor_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"recovery_conditions": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"recovery_timeframe": {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
				},
				"trigger_conditions": resourceSchemeTemplateMetricTriggerConditions(),
				"no_data_conditions": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"no_data_timeframe": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"no_data_alert_state": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"notify_no_data": {
								Type:     schema.TypeBool,
								Optional: true,
							},
						},
					},
				},
				"alarm_tags": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"auto_tags": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"custom_tags": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"custom_annotations": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSchemeTemplateMetricTriggerConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"metric_query_mode": {
					Type:     schema.TypeString,
					Required: true,
				},
				"metric_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"promql": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"aggregation_window": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"query_match": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"aggregate_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"metric_labels": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"aggregation_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"operator": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"thresholds": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"trigger_times": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"trigger_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"trigger_interval": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"expression": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"mix_promql": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"metric_statistic_method": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"metric_namespace": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"metric_unit": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"promql_expr": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"promql_for": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"aom_monitor_level": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func resourceAlarmRulesTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v4/{project_id}/alarm-rules-template"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildCreateAlarmRulesTemplateBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating alarm rules template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening creating alarm rules template response: %s", err)
	}

	id := utils.PathSearch("alarm_rule_templates[0].alarm_rule_template_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating AOM alarm rules template: can not find alarm_rule_template_id in return")
	}

	d.SetId(id)

	return resourceAlarmRulesTemplateRead(ctx, d, meta)
}

func buildCreateAlarmRulesTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_rule_template_name":        d.Get("name"),
		"alarm_rule_template_type":        d.Get("type"),
		"alarm_template_spec_list":        buildAlarmTemplateSpecList(d),
		"templating":                      buildTemplating(d),
		"alarm_rule_template_description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func buildAlarmTemplateSpecList(d *schema.ResourceData) interface{} {
	paramsList := d.Get("alarm_template_spec_list").([]interface{})
	rst := make([]map[string]interface{}, 0)
	for _, rawParams := range paramsList {
		params, ok := rawParams.(map[string]interface{})
		if ok {
			m := utils.RemoveNil(map[string]interface{}{
				"related_cloud_service":        utils.ValueIgnoreEmpty(params["related_cloud_service"]),
				"related_cce_clusters":         utils.ValueIgnoreEmpty(params["related_cce_clusters"].(*schema.Set).List()),
				"related_prometheus_instances": utils.ValueIgnoreEmpty(params["related_prometheus_instances"].(*schema.Set).List()),
				"alarm_notification":           buildTemplateAlarmNotification(params["alarm_notification"].([]interface{})),
				"alarm_template_spec_items":    buildTemplateAlarmTemplateSpecItems(params["alarm_template_spec_items"].(*schema.Set).List()),
			})
			if !reflect.DeepEqual(m, map[string]interface{}{}) {
				rst = append(rst, m)
			}
		}
	}

	return &rst
}

func buildTemplateAlarmTemplateSpecItems(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0)
	for _, rawParams := range paramsList {
		params := rawParams.(map[string]interface{})
		// remove empty rule
		if params["alarm_rule_name"] == "" {
			continue
		}
		m := map[string]interface{}{
			"alarm_template_name":        params["alarm_rule_name"],
			"alarm_template_spec_type":   params["alarm_rule_type"],
			"desc":                       utils.ValueIgnoreEmpty(params["alarm_rule_description"]),
			"event_alarm_template_spec":  buildTemplateEventAlarmTemplateSpec(params["event_alarm_spec"].([]interface{})),
			"metric_alarm_template_spec": buildTemplateMetricAlarmTemplateSpec(params["metric_alarm_spec"].([]interface{})),
		}
		rst = append(rst, m)
	}

	return rst
}

func buildTemplateEventAlarmTemplateSpec(rawParams []interface{}) map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	params, ok := rawParams[0].(map[string]interface{})
	if !ok {
		return nil
	}
	rst := map[string]interface{}{
		"alarm_subtype":            utils.ValueIgnoreEmpty(params["alarm_subtype"]),
		"alarm_source":             utils.ValueIgnoreEmpty(params["alarm_source"]),
		"event_source":             utils.ValueIgnoreEmpty(params["event_source"]),
		"trigger_conditions":       buildAlarmRuleEventAlarmSpecTriggerConditions(params["trigger_conditions"].(*schema.Set).List()),
		"monitor_objects":          utils.ValueIgnoreEmpty(params["monitor_objects"]),
		"monitor_object_templates": utils.ValueIgnoreEmpty(params["monitor_object_templates"].(*schema.Set).List()),
	}

	return rst
}

func buildTemplateMetricAlarmTemplateSpec(rawParams []interface{}) map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	params, ok := rawParams[0].(map[string]interface{})
	if !ok {
		return nil
	}
	rst := map[string]interface{}{
		"alarm_subtype":       utils.ValueIgnoreEmpty(params["alarm_subtype"]),
		"alarm_source":        utils.ValueIgnoreEmpty(params["alarm_source"]),
		"monitor_type":        utils.ValueIgnoreEmpty(params["monitor_type"]),
		"recovery_conditions": buildTemplateMetricAlarmTemplateSpecRecoveryConditions(params["recovery_conditions"].([]interface{})),
		"no_data_conditions":  buildAlarmRuleMetricAlarmSpecNoDataConditions(params["no_data_conditions"].([]interface{})),
		"alarm_tags":          buildTemplateMetricAlarmSpecAlarmTags(params["alarm_tags"].([]interface{})),
		"trigger_conditions":  buildTemplateMetricAlarmTemplateSpecTriggerConditions(params["trigger_conditions"].(*schema.Set).List()),
	}

	return rst
}

func buildTemplateMetricAlarmSpecAlarmTags(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := map[string]interface{}{}
	params, ok := paramsList[0].(map[string]interface{})
	if ok {
		rst["auto_tags"] = params["auto_tags"].(*schema.Set).List()
		rst["custom_annotations"] = params["custom_annotations"].(*schema.Set).List()
		rst["custom_tags"] = params["custom_tags"].(*schema.Set).List()
	}

	return []interface{}{rst}
}

func buildTemplateMetricAlarmTemplateSpecRecoveryConditions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := map[string]interface{}{}
	params, ok := paramsList[0].(map[string]interface{})
	if ok {
		rst["recovery_timeframe"] = utils.ValueIgnoreEmpty(params["recovery_timeframe"])
	}

	return rst
}

func buildTemplateMetricAlarmTemplateSpecTriggerConditions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0)
	for _, rawParams := range paramsList {
		params := rawParams.(map[string]interface{})
		// remove empty condition
		if params["metric_query_mode"] == "" {
			continue
		}

		m := map[string]interface{}{
			"metric_query_mode":       params["metric_query_mode"],
			"metric_name":             utils.ValueIgnoreEmpty(params["metric_name"]),
			"promql":                  utils.ValueIgnoreEmpty(params["promql"]),
			"aggregation_window":      utils.ValueIgnoreEmpty(params["aggregation_window"]),
			"query_match":             utils.ValueIgnoreEmpty(params["query_match"]),
			"aggregate_type":          utils.ValueIgnoreEmpty(params["aggregate_type"]),
			"metric_labels":           utils.ValueIgnoreEmpty(params["metric_labels"].(*schema.Set).List()),
			"aggregation_type":        utils.ValueIgnoreEmpty(params["aggregation_type"]),
			"operator":                utils.ValueIgnoreEmpty(params["operator"]),
			"thresholds":              utils.ValueIgnoreEmpty(params["thresholds"]),
			"trigger_times":           utils.ValueIgnoreEmpty(params["trigger_times"]),
			"trigger_type":            utils.ValueIgnoreEmpty(params["trigger_type"]),
			"trigger_interval":        utils.ValueIgnoreEmpty(params["trigger_interval"]),
			"expression":              utils.ValueIgnoreEmpty(params["expression"]),
			"mix_promql":              utils.ValueIgnoreEmpty(params["mix_promql"]),
			"metric_namespace":        utils.ValueIgnoreEmpty(params["metric_namespace"]),
			"metric_unit":             utils.ValueIgnoreEmpty(params["metric_unit"]),
			"promql_expr":             utils.ValueIgnoreEmpty(params["promql_expr"]),
			"promql_for":              utils.ValueIgnoreEmpty(params["promql_for"]),
			"aom_monitor_level":       utils.ValueIgnoreEmpty(params["aom_monitor_level"]),
			"metric_statistic_method": utils.ValueIgnoreEmpty(params["metric_statistic_method"]),
		}
		rst = append(rst, m)
	}

	return rst
}

func buildTemplateAlarmNotification(paramsList []interface{}) map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	params := paramsList[0].(map[string]interface{})
	rst := map[string]interface{}{
		"notification_type":         params["notification_type"],
		"route_group_enable":        utils.ValueIgnoreEmpty(params["route_group_enable"]),
		"route_group_rule":          utils.ValueIgnoreEmpty(params["route_group_rule"]),
		"notification_enable":       utils.ValueIgnoreEmpty(params["notification_enable"]),
		"bind_notification_rule_id": utils.ValueIgnoreEmpty(params["bind_notification_rule_id"]),
		"notify_resolved":           utils.ValueIgnoreEmpty(params["notify_resolved"]),
		"notify_triggered":          utils.ValueIgnoreEmpty(params["notify_triggered"]),
		"notify_frequency":          utils.ValueIgnoreEmpty(params["notify_frequency"]),
	}

	return rst
}

func buildTemplating(d *schema.ResourceData) *map[string]interface{} {
	empty := map[string]interface{}{}
	paramsList := d.Get("templating").([]interface{})
	if len(paramsList) == 0 {
		return &empty
	}
	params := paramsList[0].(map[string]interface{})
	rst := map[string]interface{}{
		"list": buildTemplatingList(params["list"].(*schema.Set).List()),
	}

	return &rst
}

func buildTemplatingList(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, rawParams := range paramsList {
		params := rawParams.(map[string]interface{})
		m := map[string]interface{}{
			"name": params["name"],
			// have to send empty value for query
			"query":       params["query"],
			"type":        utils.ValueIgnoreEmpty(params["type"]),
			"description": utils.ValueIgnoreEmpty(params["description"]),
		}
		rst = append(rst, m)
	}

	return rst
}

func resourceAlarmRulesTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	template, err := getAlarmRulesTemplate(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving alarm rules template")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("alarm_rule_template_name", template, nil)),
		d.Set("type", utils.PathSearch("alarm_rule_template_type", template, nil)),
		d.Set("description", utils.PathSearch("alarm_rule_template_description", template, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", template, float64(0)).(float64))/1000, true)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("modify_time", template, float64(0)).(float64))/1000, true)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", template, nil)),
		d.Set("alarm_template_spec_list", flattenTemplateSpecList(
			utils.PathSearch("alarm_template_spec_list", template, make([]interface{}, 0)).([]interface{}))),
		d.Set("templating", flattenTemplating(utils.PathSearch("templating", template, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAlarmRulesTemplate(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getHttpUrl := "v4/{project_id}/alarm-rules-template?id={id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":          "application/json",
			"Enterprise-Project-Id": "all_granted_eps",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving alarm rules template: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening alarm rules template: %s", err)
	}

	template := utils.PathSearch("[]|[0]", getRespBody, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return template, nil
}

func flattenTemplateSpecList(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"related_cloud_service":        utils.PathSearch("related_cloud_service", params, nil),
			"related_cce_clusters":         utils.PathSearch("related_cce_clusters", params, nil),
			"related_prometheus_instances": utils.PathSearch("related_prometheus_instances", params, nil),
			"alarm_notification":           flattenV4AlarmNotifications(utils.PathSearch("alarm_notification", params, nil)),
			"alarm_template_spec_items": flattenTemplateSpecListItems(
				utils.PathSearch("alarm_template_spec_items", params, make([]interface{}, 0)).([]interface{})),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenTemplateSpecListItems(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"alarm_rule_name":        utils.PathSearch("alarm_template_name", params, nil),
			"alarm_rule_type":        utils.PathSearch("alarm_template_spec_type", params, nil),
			"alarm_rule_description": utils.PathSearch("desc", params, nil),
			"event_alarm_spec":       flattenTemplateEventAlarmSpec(utils.PathSearch("event_alarm_template_spec", params, nil)),
			"metric_alarm_spec":      flattenTemplateMetricAlarmSpec(utils.PathSearch("metric_alarm_template_spec", params, nil)),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenTemplateEventAlarmSpec(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"alarm_subtype": utils.PathSearch("alarm_subtype", params, nil),
		"alarm_source":  utils.PathSearch("alarm_source", params, nil),
		"event_source":  utils.PathSearch("event_source", params, nil),
		"trigger_conditions": flattenV4EventTriggerConditions(
			utils.PathSearch("trigger_conditions", params, make([]interface{}, 0)).([]interface{})),
		"monitor_objects":          utils.PathSearch("monitor_objects", params, nil),
		"monitor_object_templates": utils.PathSearch("monitor_object_templates", params, nil),
	}

	return []map[string]interface{}{rst}
}

func flattenTemplateMetricAlarmSpec(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"alarm_subtype": utils.PathSearch("alarm_subtype", params, nil),
		"alarm_source":  utils.PathSearch("alarm_source", params, nil),
		"monitor_type":  utils.PathSearch("monitor_type", params, nil),
		"no_data_conditions": flattenV4NoDataConditionsMap(
			utils.PathSearch("no_data_conditions", params, make([]interface{}, 0)).([]interface{})),
		"alarm_tags": flattenV4AlarmTags(
			utils.PathSearch("alarm_tags", params, make([]interface{}, 0)).([]interface{})),
		"trigger_conditions": flattenV4TriggerConditions(
			utils.PathSearch("trigger_conditions", params, make([]interface{}, 0)).([]interface{})),
		"recovery_conditions": flattenV4MetricRecoveryConditions(utils.PathSearch("recovery_conditions", params, nil)),
	}

	return []map[string]interface{}{rst}
}

func flattenTemplating(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"list": flattenTemplatingList(utils.PathSearch("list", params, make([]interface{}, 0)).([]interface{})),
	}

	return []map[string]interface{}{rst}
}

func flattenTemplatingList(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":        utils.PathSearch("name", params, nil),
			"type":        utils.PathSearch("type", params, nil),
			"query":       utils.PathSearch("query", params, nil),
			"description": utils.PathSearch("description", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func resourceAlarmRulesTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateAlarmRuleChanges := []string{
		"type",
		"description",
		"alarm_template_spec_list",
		"templating",
	}

	if d.HasChanges(updateAlarmRuleChanges...) {
		updateHttpUrl := "v4/{project_id}/alarm-rules-template"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
			JSONBody:         utils.RemoveNil(buildUpdateAlarmRulesTemplateBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating alarm rules template: %s", err)
		}
	}

	return resourceAlarmRulesTemplateRead(ctx, d, meta)
}

func buildUpdateAlarmRulesTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_rule_template_id":          d.Id(),
		"alarm_rule_template_name":        d.Get("name"),
		"alarm_rule_template_type":        d.Get("type"),
		"alarm_template_spec_list":        buildAlarmTemplateSpecList(d),
		"templating":                      buildTemplating(d),
		"alarm_rule_template_description": d.Get("description"),
	}

	return bodyParams
}

func resourceAlarmRulesTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteHttpUrl := "v4/{project_id}/alarm-rules-template"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody: map[string]interface{}{
			"alarm_rule_templates": []string{d.Id()},
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "AOM.02018001"),
			"error deleting alarm rules template")
	}

	return nil
}
