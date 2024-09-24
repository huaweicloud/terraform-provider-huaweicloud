package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// MetricAlarmSpec 指标类告警结构。
type MetricAlarmSpec struct {

	// 监控类型。 - “all_metric”：全量指标 - “promql”：PromQL - “resource”：（日落）资源类型
	MonitorType MetricAlarmSpecMonitorType `json:"monitor_type"`

	// 无数据处理。
	NoDataConditions *[]NoDataCondition `json:"no_data_conditions,omitempty"`

	// 告警标签。
	AlarmTags []AlarmTags `json:"alarm_tags"`

	// 监控对象列表。
	MonitorObjects *[]map[string]string `json:"monitor_objects,omitempty"`

	RecoveryConditions *RecoveryCondition `json:"recovery_conditions"`

	// 触发条件。
	TriggerConditions []TriggerCondition `json:"trigger_conditions"`

	// 是否绑定告警规则模版（废弃）。
	AlarmRuleTemplateBindEnable *bool `json:"alarm_rule_template_bind_enable,omitempty"`

	// 告警规则模版id（废弃）。
	AlarmRuleTemplateId *string `json:"alarm_rule_template_id,omitempty"`
}

func (o MetricAlarmSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MetricAlarmSpec struct{}"
	}

	return strings.Join([]string{"MetricAlarmSpec", string(data)}, " ")
}

type MetricAlarmSpecMonitorType struct {
	value string
}

type MetricAlarmSpecMonitorTypeEnum struct {
	ALL_METRIC MetricAlarmSpecMonitorType
	PROMQL     MetricAlarmSpecMonitorType
	RESOURCE   MetricAlarmSpecMonitorType
}

func GetMetricAlarmSpecMonitorTypeEnum() MetricAlarmSpecMonitorTypeEnum {
	return MetricAlarmSpecMonitorTypeEnum{
		ALL_METRIC: MetricAlarmSpecMonitorType{
			value: "all_metric",
		},
		PROMQL: MetricAlarmSpecMonitorType{
			value: "promql",
		},
		RESOURCE: MetricAlarmSpecMonitorType{
			value: "resource",
		},
	}
}

func (c MetricAlarmSpecMonitorType) Value() string {
	return c.value
}

func (c MetricAlarmSpecMonitorType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MetricAlarmSpecMonitorType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
