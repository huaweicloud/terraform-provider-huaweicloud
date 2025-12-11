package aom

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v4/{project_id}/alarm-rules
// @API AOM GET /v4/{project_id}/alarm-rules
// @API AOM DELETE /v4/{project_id}/alarm-rules
func ResourceAlarmRuleV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmRuleV4Create,
		ReadContext:   resourceAlarmRuleV4Read,
		UpdateContext: resourceAlarmRuleV4Update,
		DeleteContext: resourceAlarmRuleV4Delete,

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
				ForceNew: true,
			},
			"alarm_notifications": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     resourceSchemeV4AlarmNotification(),
			},
			"event_alarm_spec": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"metric_alarm_spec"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_source": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"event_source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"trigger_conditions": resourceSchemeV4EventTriggerConditions(),
						"monitor_objects": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
								Elem: &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
			"metric_alarm_spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"recovery_conditions": {
							Type:     schema.TypeList,
							Required: true,
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
						"trigger_conditions": resourceSchemeV4MetricTriggerConditions(),
						"no_data_conditions": resourceSchemeV4NoDataConditions(),
						"alarm_tags":         resourceSchemeV4AlarmTags(),
						"monitor_objects": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
								Elem: &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
			"prom_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alarm_rule_id": {
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
	}
}

func resourceSchemeV4AlarmNotification() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"notification_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route_group_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"route_group_rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"bind_notification_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_resolved": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"notify_triggered": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"notify_frequency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSchemeV4EventTriggerConditions() *schema.Schema {
	return &schema.Schema{
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
					Computed: true,
				},
			},
		},
		Set: resourceEventTriggerConditionHash,
	}
}

func resourceEventTriggerConditionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["trigger_type"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["trigger_type"].(string)))
	}
	if m["event_name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["event_name"].(string)))
	}
	if m["thresholds"] != nil {
		buf.WriteString(fmt.Sprintf("%v-", m["thresholds"]))
	}
	if m["aggregation_window"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["aggregation_window"].(int)))
	}
	if m["operator"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["operator"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceSchemeV4NoDataConditions() *schema.Schema {
	return &schema.Schema{
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
					Computed: true,
				},
			},
		},
	}
}

func resourceSchemeV4AlarmTags() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
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
	}
}

func resourceSchemeV4MetricTriggerConditions() *schema.Schema {
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
					Required: true,
				},
				"promql": {
					Type:     schema.TypeString,
					Required: true,
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
				"query_param": {
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
					Computed: true,
				},
				"aom_monitor_level": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
		Set: resourceMetricTriggerConditionHash,
	}
}

func normalizeJsonStringForSet(jsonStr string) string {
	if jsonStr == "" {
		return jsonStr
	}
	normalized, err := utils.NormalizeJsonString(jsonStr)
	if err != nil {
		// If normalization fails, return original string
		return jsonStr
	}
	return normalized
}

func resourceMetricTriggerConditionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["metric_query_mode"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["metric_query_mode"].(string)))
	}
	if m["metric_name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["metric_name"].(string)))
	}
	if m["promql"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["promql"].(string)))
	}
	if m["aggregation_window"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["aggregation_window"].(string)))
	}
	if m["query_match"] != nil {
		queryMatchStr := ""
		if str, ok := m["query_match"].(string); ok && str != "" {
			queryMatchStr = normalizeJsonStringForSet(str)
		}
		buf.WriteString(fmt.Sprintf("%s-", queryMatchStr))
	}
	if m["aggregate_type"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["aggregate_type"].(string)))
	}
	if m["aggregation_type"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["aggregation_type"].(string)))
	}
	if m["operator"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["operator"].(string)))
	}
	if m["thresholds"] != nil {
		buf.WriteString(fmt.Sprintf("%v-", m["thresholds"]))
	}
	if m["trigger_times"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["trigger_times"].(string)))
	}
	if m["trigger_type"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["trigger_type"].(string)))
	}
	if m["expression"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["expression"].(string)))
	}
	if m["trigger_interval"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["trigger_interval"].(string)))
	}
	if m["mix_promql"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["mix_promql"].(string)))
	}
	if m["metric_statistic_method"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["metric_statistic_method"].(string)))
	}
	if m["metric_labels"] != nil {
		buf.WriteString(fmt.Sprintf("%v-", m["metric_labels"]))
	}
	if m["query_param"] != nil {
		queryParamStr := ""
		if str, ok := m["query_param"].(string); ok && str != "" {
			queryParamStr = normalizeJsonStringForSet(str)
		}
		buf.WriteString(fmt.Sprintf("%s-", queryParamStr))
	}
	if m["metric_namespace"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["metric_namespace"].(string)))
	}
	if m["metric_unit"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["metric_unit"].(string)))
	}
	if m["promql_expr"] != nil {
		buf.WriteString(fmt.Sprintf("%v-", m["promql_expr"]))
	}
	if m["aom_monitor_level"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["aom_monitor_level"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceAlarmRuleV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v4/{project_id}/alarm-rules?action_id=add-alarm-action"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildAlarmRuleBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating the alarm rule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceAlarmRuleV4Read(ctx, d, meta)
}

func buildAlarmRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_rule_name":        d.Get("name"),
		"alarm_rule_type":        d.Get("type"),
		"alarm_notifications":    buildAlarmRuleAlarmNotifications(d),
		"metric_alarm_spec":      buildAlarmRuleMetricAlarmSpec(d),
		"event_alarm_spec":       buildAlarmRuleEventAlarmSpec(d),
		"prom_instance_id":       utils.ValueIgnoreEmpty(d.Get("prom_instance_id")),
		"alarm_rule_enable":      utils.ValueIgnoreEmpty(d.Get("enable")),
		"alarm_rule_description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func buildAlarmRuleAlarmNotifications(d *schema.ResourceData) map[string]interface{} {
	params := d.Get("alarm_notifications.0").(map[string]interface{})
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

func buildAlarmRuleEventAlarmSpec(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("event_alarm_spec").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := rawParams[0].(map[string]interface{})
	rst := map[string]interface{}{
		"alarm_source":       params["alarm_source"],
		"event_source":       params["event_source"],
		"trigger_conditions": buildAlarmRuleEventAlarmSpecTriggerConditions(params["trigger_conditions"].(*schema.Set).List()),
		"monitor_objects":    utils.ValueIgnoreEmpty(params["monitor_objects"]),
	}

	return rst
}

func buildAlarmRuleEventAlarmSpecTriggerConditions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0)
	for _, rawParams := range paramsList {
		params := rawParams.(map[string]interface{})
		// remove empty condition
		if params["trigger_type"] == "" {
			continue
		}
		m := map[string]interface{}{
			"trigger_type":       params["trigger_type"],
			"thresholds":         utils.ValueIgnoreEmpty(params["thresholds"]),
			"event_name":         utils.ValueIgnoreEmpty(params["event_name"]),
			"aggregation_window": utils.ValueIgnoreEmpty(params["aggregation_window"]),
			"operator":           utils.ValueIgnoreEmpty(params["operator"]),
			"frequency":          utils.ValueIgnoreEmpty(params["frequency"]),
		}
		rst = append(rst, m)
	}

	return rst
}

func buildAlarmRuleMetricAlarmSpec(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("metric_alarm_spec").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := rawParams[0].(map[string]interface{})
	rst := map[string]interface{}{
		"monitor_type":        params["monitor_type"],
		"recovery_conditions": buildAlarmRuleMetricAlarmSpecRecoveryConditions(params["recovery_conditions"].([]interface{})),
		"no_data_conditions":  buildAlarmRuleMetricAlarmSpecNoDataConditions(params["no_data_conditions"].([]interface{})),
		"alarm_tags":          buildAlarmRuleMetricAlarmSpecAlarmTags(params["alarm_tags"].([]interface{})),
		"trigger_conditions":  buildAlarmRuleMetricAlarmSpecTriggerConditions(params["trigger_conditions"].(*schema.Set).List()),
		"monitor_objects":     utils.ValueIgnoreEmpty(params["monitor_objects"]),
	}

	return rst
}

func buildAlarmRuleMetricAlarmSpecRecoveryConditions(paramsList []interface{}) interface{} {
	rst := map[string]interface{}{}
	params, ok := paramsList[0].(map[string]interface{})
	// return empty structure for `recovery_conditions` is required
	if ok {
		rst["recovery_timeframe"] = utils.ValueIgnoreEmpty(params["recovery_timeframe"])
	}

	return &rst
}

func buildAlarmRuleMetricAlarmSpecNoDataConditions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	var rst []map[string]interface{}
	params, ok := paramsList[0].(map[string]interface{})
	if ok {
		m := map[string]interface{}{
			"notify_no_data":      utils.ValueIgnoreEmpty(params["notify_no_data"]),
			"no_data_timeframe":   utils.ValueIgnoreEmpty(params["no_data_timeframe"]),
			"no_data_alert_state": utils.ValueIgnoreEmpty(params["no_data_alert_state"]),
		}
		rst = append(rst, m)
	}

	return rst
}

func buildAlarmRuleMetricAlarmSpecAlarmTags(paramsList []interface{}) interface{} {
	// have to send empty structure
	rst := map[string]interface{}{
		"auto_tags":          []interface{}{},
		"custom_annotations": []interface{}{},
		"custom_tags":        []interface{}{},
	}
	if len(paramsList) != 0 {
		params, ok := paramsList[0].(map[string]interface{})
		if ok {
			rst["auto_tags"] = params["auto_tags"].(*schema.Set).List()
			rst["custom_annotations"] = params["custom_annotations"].(*schema.Set).List()
			rst["custom_tags"] = params["custom_tags"].(*schema.Set).List()
		}
	}

	return []map[string]interface{}{rst}
}

func buildAlarmRuleMetricAlarmSpecTriggerConditions(paramsList []interface{}) interface{} {
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
			"metric_name":             params["metric_name"],
			"promql":                  params["promql"],
			"aggregation_window":      utils.ValueIgnoreEmpty(params["aggregation_window"]),
			"query_match":             utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("query_match", params, "").(string))),
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
			"query_param":             utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("query_param", params, "").(string))),
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

func getAlarmRule(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getHttpUrl := "v4/{project_id}/alarm-rules?name={name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{name}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the alarm rule: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening the alarm rule: %s", err)
	}

	rule := utils.PathSearch("alarm_rules|[0]", getRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func resourceAlarmRuleV4Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	rule, err := getAlarmRule(cfg, client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the alarm rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("alarm_rule_name", rule, nil)),
		d.Set("type", utils.PathSearch("alarm_rule_type", rule, nil)),
		d.Set("alarm_notifications", flattenV4AlarmNotifications(utils.PathSearch("alarm_notifications", rule, nil))),
		d.Set("metric_alarm_spec", flattenV4MetricAlarmSpec(utils.PathSearch("metric_alarm_spec", rule, nil))),
		d.Set("event_alarm_spec", flattenV4EventAlarmSpec(utils.PathSearch("event_alarm_spec", rule, nil))),
		d.Set("prom_instance_id", utils.PathSearch("prom_instance_id", rule, nil)),
		d.Set("enable", utils.PathSearch("alarm_rule_enable", rule, nil)),
		d.Set("description", utils.PathSearch("alarm_rule_description", rule, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("alarm_create_time", rule, float64(0)).(float64))/1000, true)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("alarm_update_time", rule, float64(0)).(float64))/1000, true)),
		d.Set("status", utils.PathSearch("alarm_rule_status", rule, nil)),
		d.Set("alarm_rule_id", strconv.FormatFloat(utils.PathSearch("alarm_rule_id", rule, float64(0)).(float64), 'f', -1, 64)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", rule, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV4AlarmNotifications(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"notification_type":         utils.PathSearch("notification_type", params, nil),
		"route_group_enable":        utils.PathSearch("route_group_enable", params, nil),
		"route_group_rule":          utils.PathSearch("route_group_rule", params, nil),
		"notification_enable":       utils.PathSearch("notification_enable", params, nil),
		"bind_notification_rule_id": utils.PathSearch("bind_notification_rule_id", params, nil),
		"notify_resolved":           utils.PathSearch("notify_resolved", params, nil),
		"notify_triggered":          utils.PathSearch("notify_triggered", params, nil),
		"notify_frequency":          strconv.FormatFloat(utils.PathSearch("notify_frequency", params, float64(0)).(float64), 'f', -1, 64),
	}

	return []map[string]interface{}{rst}
}

func flattenV4EventAlarmSpec(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"alarm_source": utils.PathSearch("alarm_source", params, nil),
		"event_source": utils.PathSearch("event_source", params, nil),
		"trigger_conditions": flattenV4EventTriggerConditions(
			utils.PathSearch("trigger_conditions", params, make([]interface{}, 0)).([]interface{})),
		"monitor_objects": utils.PathSearch("monitor_objects", params, nil),
	}

	return []map[string]interface{}{rst}
}

func flattenV4EventTriggerConditions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"event_name":         utils.PathSearch("event_name", params, nil),
			"trigger_type":       utils.PathSearch("trigger_type", params, nil),
			"aggregation_window": utils.PathSearch("aggregation_window", params, nil),
			"operator":           utils.PathSearch("operator", params, nil),
			"thresholds":         utils.PathSearch("thresholds", params, nil),
			"frequency":          utils.PathSearch("frequency", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenV4MetricAlarmSpec(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"monitor_type": utils.PathSearch("monitor_type", params, nil),
		"no_data_conditions": flattenV4NoDataConditionsMap(
			utils.PathSearch("no_data_conditions", params, make([]interface{}, 0)).([]interface{})),
		"alarm_tags": flattenV4AlarmTags(
			utils.PathSearch("alarm_tags", params, make([]interface{}, 0)).([]interface{})),
		"trigger_conditions": flattenV4TriggerConditions(
			utils.PathSearch("trigger_conditions", params, make([]interface{}, 0)).([]interface{})),
		"recovery_conditions": flattenV4MetricRecoveryConditions(utils.PathSearch("recovery_conditions", params, nil)),
		"monitor_objects":     utils.PathSearch("monitor_objects", params, nil),
	}

	return []map[string]interface{}{rst}
}

func flattenV4NoDataConditionsMap(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"no_data_timeframe":   utils.PathSearch("no_data_timeframe", params, nil),
			"no_data_alert_state": utils.PathSearch("no_data_alert_state", params, nil),
			"notify_no_data":      utils.PathSearch("notify_no_data", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenV4AlarmTags(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"auto_tags":          utils.PathSearch("auto_tags", params, nil),
			"custom_tags":        utils.PathSearch("custom_tags", params, nil),
			"custom_annotations": utils.PathSearch("custom_annotations", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenV4MetricRecoveryConditions(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"recovery_timeframe": utils.PathSearch("recovery_timeframe", params, nil),
	}

	return []map[string]interface{}{rst}
}

func flattenV4TriggerConditions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		// query_param is unset, it's a structure in return
		queryMatchStr := ""
		if queryMatch := utils.PathSearch("query_match", params, ""); queryMatch != nil {
			if str, ok := queryMatch.(string); ok {
				queryMatchStr = normalizeJsonStringForSet(str)
			}
		}

		queryParamStr := ""
		if queryParam := utils.PathSearch("query_param", params, nil); queryParam != nil {
			queryParamStr = normalizeJsonStringForSet(utils.JsonToString(queryParam))
		}

		m := map[string]interface{}{
			"metric_query_mode":  utils.PathSearch("metric_query_mode", params, nil),
			"metric_name":        utils.PathSearch("metric_name", params, nil),
			"metric_labels":      utils.PathSearch("metric_labels", params, nil),
			"promql":             utils.PathSearch("promql", params, nil),
			"trigger_times":      strconv.FormatFloat(utils.PathSearch("trigger_times", params, float64(0)).(float64), 'f', -1, 64),
			"trigger_interval":   utils.PathSearch("trigger_interval", params, nil),
			"trigger_type":       utils.PathSearch("trigger_type", params, nil),
			"aggregation_type":   utils.PathSearch("aggregation_type", params, nil),
			"aggregation_window": utils.PathSearch("aggregation_window", params, nil),
			"operator":           utils.PathSearch("operator", params, nil),
			"thresholds":         utils.PathSearch("thresholds", params, nil),
			// Normalize JSON strings to ensure consistent field order for TypeSet hash calculation
			"query_match":             queryMatchStr,
			"query_param":             queryParamStr,
			"aggregate_type":          utils.PathSearch("aggregate_type", params, nil),
			"metric_statistic_method": utils.PathSearch("metric_statistic_method", params, nil),
			"expression":              utils.PathSearch("expression", params, nil),
			"mix_promql":              utils.PathSearch("mix_promql", params, nil),
			"metric_namespace":        utils.PathSearch("metric_namespace", params, nil),
			"metric_unit":             utils.PathSearch("metric_unit", params, nil),
			"promql_expr":             utils.PathSearch("promql_expr", params, nil),
			"promql_for":              utils.PathSearch("promql_for", params, nil),
			"aom_monitor_level":       utils.PathSearch("aom_monitor_level", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func resourceAlarmRuleV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateAlarmRuleChanges := []string{
		"description",
		"enable",
		"prom_instance_id",
		"alarm_notifications",
		"event_alarm_spec",
		"metric_alarm_spec",
	}

	if d.HasChanges(updateAlarmRuleChanges...) {
		updateHttpUrl := "v4/{project_id}/alarm-rules?action_id=update-alarm-action"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
			JSONBody:         utils.RemoveNil(buildAlarmRuleBodyParams(d)),
		}

		_, err = client.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating the alarm rule: %s", err)
		}
	}

	return resourceAlarmRuleV4Read(ctx, d, meta)
}

func resourceAlarmRuleV4Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	// DELETE will return success even deleting a non exist rule
	_, err = getAlarmRule(cfg, client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting the alarm rule")
	}

	deleteHttpUrl := "v4/{project_id}/alarm-rules"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody: map[string]interface{}{
			"alarm_rules": []string{d.Id()},
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting the alarm rule")
	}

	return nil
}
