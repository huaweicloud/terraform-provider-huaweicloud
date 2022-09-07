package aom

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io"
	"time"
)

func ResourceAlarmPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAlarmPolicyRead,
		CreateContext: resourceAlarmPolicyCreate,
		DeleteContext: resourceAlarmPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"alarm_rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alarm_rule_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"alarm_rule_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"alarm_rule_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"alarm_rule_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric_alarm_spec":   schemaMetricAlarmSpe(),
			"event_alarm_spec":    schemeEventAlarmSpec(),
			"alarm_notifications": schemeAlarmNotifications(),
		},
	}
}

func schemaMetricAlarmSpe() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"monitor_type": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"resource_kind": {
					Type:     schema.TypeString,
					ForceNew: true,
					Optional: true,
					Computed: true,
				},
				"metric_kind": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"alarm_rule_template_bind_enable": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"alarm_rule_template_id": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"no_data_conditions": schemeNoDataConditions(),
				"alarm_tags":         schemeAlarmTags(),
				"trigger_conditions": schemeTriggerConditions(),
				"monitor_objects":    schemeMonitorObjects(),
				"recovery_conditions": {
					Type:     schema.TypeMap,
					Optional: true,
					ForceNew: true,
					Computed: true,
					Elem:     schema.TypeInt,
				},
			},
		},
	}
}

func schemeEventAlarmSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"event_source": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"monitor_objects":    schemeMonitorObjects(),
				"no_data_conditions": schemeNoDataConditions(),
				"alarm_tags":         schemeAlarmTags(),
			},
		},
	}
}

func schemeNoDataConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"no_data_timeframe": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"no_data_alert_state": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"notify_no_data": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},
		},
	}

}

func schemeAlarmTags() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_tags": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"custom_tags": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"custom_annotations": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func schemeTriggerConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"metric_query_mode": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"metric_namespace": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"metric_name": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"metric_unit": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"metric_labels": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"promql": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"trigger_times": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"trigger_interval": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"aggregation_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"aggregation_window": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"operator": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"trigger_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"thresholds": {
					Type:     schema.TypeMap,
					Optional: true,
					ForceNew: true,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func schemeMonitorObjects() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString}},
	}
}

func schemeAlarmNotifications() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"notification_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"route_group_enable": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"route_group_rule": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"inhibit_enable": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"inhibit_rule": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"notification_enable": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"bind_notification_rule_id": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
				"notify_resolved": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},
			},
		},
	}
}

func resourceAlarmPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := httpclient_go.NewHttpClientGo(conf, "aom", conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	createOpts := entity.AddAlarmRuleParams{
		AlarmRuleName:        d.Get("alarm_rule_name").(string),
		EnterpriseProjectId:  conf.EnterpriseProjectID,
		AlarmRuleDescription: d.Get("alarm_rule_description").(string),
		AlarmRuleEnable:      d.Get("alarm_rule_enable").(bool),
		AlarmRuleStatus:      d.Get("alarm_rule_status").(string),
		AlarmRuleType:        d.Get("alarm_rule_type").(string),
		MetricAlarmSpec:      buildMetricAlarmSpec(d.Get("metric_alarm_spec").([]interface{})),
		EventAlarmSpec:       buildEventAlarmSpec(d.Get("event_alarm_spec").([]interface{})),
		AlarmNotifications:   buildAlarmNotifications(d.Get("alarm_notifications").([]interface{})),
	}
	region := conf.GetRegion(d)
	client.WithMethod(httpclient_go.MethodPost).WithUrl("v4/" + conf.GetProjectID(region) + "/alarm-rules?action_id=add-alarm-action").WithBody(createOpts)
	response, err := client.Do()

	if err != nil {
		return diag.Errorf("error creating AOM alarm rule %s: %s", createOpts.AlarmRuleName, err)
	}

	mErr := &multierror.Error{}
	defer response.Body.Close()

	_, err = io.ReadAll(response.Body)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error getting AOM prometheus instance fields: %s", err)
	}
	d.SetId(createOpts.AlarmRuleName)
	return resourceAlarmPolicyRead(context.TODO(), d, meta)
}

func resourceAlarmPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(conf, "aom", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	client.WithMethod(httpclient_go.MethodGet).WithUrl("v4/" + conf.GetProjectID(region) + "/alarm-rules")

	resp, err := client.Do()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM alarm rule")
	}

	mErr := &multierror.Error{}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	rlt := make([]entity.AddAlarmRuleParams, 0)

	err = json.Unmarshal(body, &rlt)

	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	for _, params := range rlt {
		if params.AlarmRuleName == d.Id() {
			d.SetId(params.AlarmRuleName)
			mErr = multierror.Append(mErr,
				d.Set("alarm_rule_status", params.AlarmRuleStatus),
				d.Set("region", conf.GetRegion(d)),
				d.Set("alarm_rule_description", params.AlarmRuleDescription),
				d.Set("alarm_rule_enable", params.AlarmRuleEnable),
				d.Set("alarm_rule_name", params.AlarmRuleName),
				d.Set("alarm_rule_type", params.AlarmRuleType),
				d.Set("alarm_notifications", buildAlarmNotificationsMap(params.AlarmNotifications)),
				d.Set("metric_alarm_spec", buildMetricAlarmSpecMap(params.MetricAlarmSpec)),
			)
			if err = mErr.ErrorOrNil(); err != nil {
				return diag.Errorf("error getting AOM alarm policy fields: %s", err)
			}
			return nil
		}
	}
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error getting AOM alarm policy fields: %s", err)
	}
	return nil
}

func buildMetricAlarmSpecMap(spec entity.MetricAlarmSpec) []map[string]interface{} {
	var m = make(map[string]interface{})
	m["monitor_type"] = spec.MonitorType
	m["resource_kind"] = spec.ResourceKind
	m["metric_kind"] = spec.MetricKind
	m["alarm_rule_template_bind_enable"] = spec.AlarmRuletemplateBindEnable
	m["alarm_rule_template_id"] = spec.AlarmRuletemplateId
	m["no_data_conditions"] = buildNoDataConditionsMap(spec.NoDataConditions)
	m["alarm_tags"] = spec.AlarmTags
	m["trigger_conditions"] = buildTriggerConditionsMap(spec.TriggerConditions)
	m["monitor_objects"] = spec.MonitorObjects
	m["recovery_conditions"] = spec.RecoveryConditions
	return []map[string]interface{}{m}
}

func buildTriggerConditionsMap(conditions []entity.TriggerCondition) interface{} {
	var ret []map[string]interface{}
	for _, condition := range conditions {
		var m = make(map[string]interface{})
		m["metric_query_mode"] = condition.MetricQueryMode
		m["metric_namespace"] = condition.MetricNamespace
		m["metric_name"] = condition.MetricName
		m["metric_labels"] = condition.MetricLabels
		m["metric_unit"] = condition.MetricUnit
		m["promql"] = condition.Promql
		m["trigger_times"] = condition.TriggerTimes
		m["trigger_interval"] = condition.TriggerInterval
		m["trigger_type"] = condition.TriggerType
		m["aggregation_window"] = condition.AggregationWindow
		m["aggregation_type"] = condition.AggregationType

		m["operator"] = condition.Operator
		m["thresholds"] = condition.Thresholds
		ret = append(ret, m)
	}
	return ret
}

func buildNoDataConditionsMap(conditions []entity.NoDataCondition) []map[string]interface{} {
	var ret []map[string]interface{}
	for _, condition := range conditions {
		var m = make(map[string]interface{})
		m["no_data_timeframe"] = condition.NoDataTimeframe
		m["no_data_alert_state"] = condition.NoDataAlertState
		m["notify_no_data"] = condition.NotifyNoData
		ret = append(ret, m)
	}
	return ret
}

func buildAlarmNotificationsMap(notifications entity.AlarmNotifications) []map[string]interface{} {
	var m = make(map[string]interface{})
	m["notification_type"] = notifications.NotificationType
	m["route_group_enable"] = notifications.RouteGroupEnable
	m["route_group_rule"] = notifications.RouteGroupRule
	m["inhibit_enable"] = notifications.InhibitEnable
	m["inhibit_rule"] = notifications.InhibitRule
	m["notification_type"] = notifications.NotificationType
	m["notify_resolved"] = notifications.NotiFyResolved
	return []map[string]interface{}{m}
}

func resourceAlarmPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := httpclient_go.NewHttpClientGo(conf, "aom", region)
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}
	client.WithMethod(httpclient_go.MethodDelete).WithUrl("v4/" + conf.GetProjectID(region) + "/alarm-rules").WithBody([]string{d.Id()})
	resp, err := client.Do()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM alarm rule")
	}

	mErr := &multierror.Error{}
	if resp.StatusCode != 200 {
		mErr = multierror.Append(mErr, fmt.Errorf("delete alarm policy failed error code: %d", resp.StatusCode))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error delete AOM alarm policy fields: %s", err)
	}

	return nil
}

func buildMetricAlarmSpec(raw interface{}) entity.MetricAlarmSpec {
	mas := make([]entity.MetricAlarmSpec, 0)
	b, err := json.Marshal(raw)
	if err != nil {
		return entity.MetricAlarmSpec{}
	}
	err = json.Unmarshal(b, &mas)
	if err != nil {
		return entity.MetricAlarmSpec{}
	}
	if len(mas) == 0 {
		return entity.MetricAlarmSpec{}
	}
	return mas[0]
}

func buildEventAlarmSpec(raw interface{}) entity.EventAlarmSpec {
	mas := make([]entity.EventAlarmSpec, 0)
	b, err := json.Marshal(raw)
	if err != nil {
		return entity.EventAlarmSpec{}
	}
	err = json.Unmarshal(b, &mas)
	if err != nil {
		return entity.EventAlarmSpec{}
	}
	if len(mas) == 0 {
		return entity.EventAlarmSpec{}
	}
	return mas[0]
}

func buildAlarmNotifications(raw interface{}) entity.AlarmNotifications {
	mas := make([]entity.AlarmNotifications, 0)
	b, err := json.Marshal(raw)
	if err != nil {
		return entity.AlarmNotifications{}
	}
	err = json.Unmarshal(b, &mas)
	if err != nil {
		return entity.AlarmNotifications{}
	}
	if len(mas) == 0 {
		return entity.AlarmNotifications{}
	}
	return mas[0]
}
