package entity

type AddAlarmRuleParams struct {
	AlarmCreateTime      int64              `json:"alarm_create_time"`
	AlarmRuleName        string             `json:"alarm_rule_name"`
	EnterpriseProjectId  string             `json:"enterprise_project_id"`
	AlarmRuleDescription string             `json:"alarm_rule_description"`
	AlarmRuleEnable      bool               `json:"alarm_rule_enable"`
	AlarmRuleStatus      string             `json:"alarm_rule_status"`
	AlarmRuleType        string             `json:"alarm_rule_type"`
	MetricAlarmSpec      MetricAlarmSpec    `json:"metric_alarm_spec"`
	EventAlarmSpec       EventAlarmSpec     `json:"event_alarm_spec"`
	AlarmNotifications   AlarmNotifications `json:"alarm_notifications"`
}

type MetricAlarmSpec struct {
	MonitorType                 string                   `json:"monitor_type"`
	ResourceKind                string                   `json:"resource_kind"`
	MetricKind                  string                   `json:"metric_kind"`
	AlarmRuletemplateBindEnable bool                     `json:"alarm_rule_template_bind_enable"`
	AlarmRuletemplateId         string                   `json:"alarm_rule_template_id"`
	NoDataConditions            []NoDataCondition        `json:"no_data_conditions"`
	AlarmTags                   []AlarmTag               `json:"alarm_tags"`
	TriggerConditions           []TriggerCondition       `json:"trigger_conditions"`
	MonitorObjects              []map[string]interface{} `json:"monitor_objects"`
	RecoveryConditions          map[string]interface{}   `json:"recovery_conditions"`
}

type EventAlarmSpec struct {
	EventSource       string                   `json:"event_source"`
	MonitorObjects    []map[string]interface{} `json:"monitor_objects"`
	TriggerConditions []TriggerCondition       `json:"trigger_conditions"`
	NoDataConditions  []NoDataCondition        `json:"no_data_conditions"`
	AlarmTags         []AlarmTag               `json:"alarm_tags"`
}

type AlarmNotifications struct {
	NotificationType       string `json:"notification_type"`
	RouteGroupEnable       bool   `json:"route_group_enable"`
	RouteGroupRule         string `json:"route_group_rule"`
	InhibitEnable          bool   `json:"inhibit_enable"`
	InhibitRule            string `json:"inhibit_rule"`
	NotiFicationEnable     bool   `json:"notification_enable"`
	BindNotificationRuleId string `json:"bind_notification_rule_id"`
	NotiFyResolved         bool   `json:"notify_resolved"`
}

type NoDataCondition struct {
	NoDataTimeframe  int    `json:"no_data_timeframe"`
	NoDataAlertState string `json:"no_data_alert_state"`
	NotifyNoData     bool   `json:"notify_no_data"`
}

type AlarmTag struct {
	AutoTags          []string `json:"auto_tags"`
	CustomTags        []string `json:"custom_tags"`
	CustomAnnotations []string `json:"custom_annotations"`
}

type TriggerCondition struct {
	MetricQueryMode   string                 `json:"metric_query_mode"`
	MetricNamespace   string                 `json:"metric_namespace"`
	MetricName        string                 `json:"metric_name"`
	MetricLabels      []string               `json:"metric_labels"`
	MetricUnit        string                 `json:"metric_unit"`
	Promql            string                 `json:"promql"`
	TriggerTimes      int                    `json:"trigger_times"`
	TriggerInterval   string                 `json:"trigger_interval"`
	TriggerType       string                 `json:"trigger_type"`
	AggregationType   string                 `json:"aggregation_type"`
	AggregationWindow interface{}            `json:"aggregation_window"`
	Operator          string                 `json:"operator"`
	Thresholds        map[string]interface{} `json:"thresholds"`
}
