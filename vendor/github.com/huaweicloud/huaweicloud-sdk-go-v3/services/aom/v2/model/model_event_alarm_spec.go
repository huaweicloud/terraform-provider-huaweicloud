package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// EventAlarmSpec 事件类告警结构。
type EventAlarmSpec struct {

	// 告警规则来源。 - “systemEvent”：系统事件 - “customEvent”：自定义事件
	AlarmSource *EventAlarmSpecAlarmSource `json:"alarm_source,omitempty"`

	// 告警来源。 - “RDS” - “EVS” - “CCE” - “LTS” - “AOM”
	EventSource *string `json:"event_source,omitempty"`

	// 监控对象列表。键值对形式，键值为： - “event_type”：通知类型 - “event_severity”：告警级别 - “event_name”：事件名称 - “namespace”：命名空间 - “clusterId”：集群id - “customField”：用户自定义字段
	MonitorObjects *[]map[string]string `json:"monitor_objects,omitempty"`

	// 触发条件。
	TriggerConditions *[]EventTriggerCondition `json:"trigger_conditions,omitempty"`

	// 是否绑定告警规则模版（废弃）。
	AlarmRuleTemplateBindEnable *bool `json:"alarm_rule_template_bind_enable,omitempty"`

	// 告警规则模版id（废弃）。
	AlarmRuleTemplateId *string `json:"alarm_rule_template_id,omitempty"`
}

func (o EventAlarmSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventAlarmSpec struct{}"
	}

	return strings.Join([]string{"EventAlarmSpec", string(data)}, " ")
}

type EventAlarmSpecAlarmSource struct {
	value string
}

type EventAlarmSpecAlarmSourceEnum struct {
	SYSTEM_EVENT EventAlarmSpecAlarmSource
	CUSTOM_EVENT EventAlarmSpecAlarmSource
}

func GetEventAlarmSpecAlarmSourceEnum() EventAlarmSpecAlarmSourceEnum {
	return EventAlarmSpecAlarmSourceEnum{
		SYSTEM_EVENT: EventAlarmSpecAlarmSource{
			value: "systemEvent",
		},
		CUSTOM_EVENT: EventAlarmSpecAlarmSource{
			value: "customEvent",
		},
	}
}

func (c EventAlarmSpecAlarmSource) Value() string {
	return c.value
}

func (c EventAlarmSpecAlarmSource) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EventAlarmSpecAlarmSource) UnmarshalJSON(b []byte) error {
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
