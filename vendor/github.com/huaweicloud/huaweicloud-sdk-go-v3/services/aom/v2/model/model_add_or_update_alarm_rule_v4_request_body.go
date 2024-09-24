package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// AddOrUpdateAlarmRuleV4RequestBody 新增或修改告警规则请求体。
type AddOrUpdateAlarmRuleV4RequestBody struct {
	AlarmNotifications *AlarmNotification `json:"alarm_notifications,omitempty"`

	// 告警规则描述。
	AlarmRuleDescription *string `json:"alarm_rule_description,omitempty"`

	// 是否启用。
	AlarmRuleEnable *bool `json:"alarm_rule_enable,omitempty"`

	// 告警规则名称。
	AlarmRuleName string `json:"alarm_rule_name"`

	// 告警规则类型。 - “metric”：指标告警规则 - “event”：事件告警规则
	AlarmRuleType AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType `json:"alarm_rule_type"`

	EventAlarmSpec *EventAlarmSpec `json:"event_alarm_spec,omitempty"`

	MetricAlarmSpec *MetricAlarmSpec `json:"metric_alarm_spec,omitempty"`

	// Prometheus实例id。
	PromInstanceId *string `json:"prom_instance_id,omitempty"`
}

func (o AddOrUpdateAlarmRuleV4RequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddOrUpdateAlarmRuleV4RequestBody struct{}"
	}

	return strings.Join([]string{"AddOrUpdateAlarmRuleV4RequestBody", string(data)}, " ")
}

type AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType struct {
	value string
}

type AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleTypeEnum struct {
	METRIC AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType
	EVENT  AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType
}

func GetAddOrUpdateAlarmRuleV4RequestBodyAlarmRuleTypeEnum() AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleTypeEnum {
	return AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleTypeEnum{
		METRIC: AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType{
			value: "metric",
		},
		EVENT: AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType{
			value: "event",
		},
	}
}

func (c AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType) Value() string {
	return c.value
}

func (c AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AddOrUpdateAlarmRuleV4RequestBodyAlarmRuleType) UnmarshalJSON(b []byte) error {
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
