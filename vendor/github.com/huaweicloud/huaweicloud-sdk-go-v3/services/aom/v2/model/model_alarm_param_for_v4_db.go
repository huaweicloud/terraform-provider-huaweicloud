package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type AlarmParamForV4Db struct {

	// 告警规则创建时间。
	AlarmCreateTime int64 `json:"alarm_create_time"`

	// 告警规则修改时间。
	AlarmUpdateTime int64 `json:"alarm_update_time"`

	// 告警规则名称。
	AlarmRuleName string `json:"alarm_rule_name"`

	// 告警规则id。
	AlarmRuleId int64 `json:"alarm_rule_id"`

	// 企业项目id。
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// Prometheus实例id。
	PromInstanceId *string `json:"prom_instance_id,omitempty"`

	// 告警规则描述。
	AlarmRuleDescription string `json:"alarm_rule_description"`

	// 是否启用。
	AlarmRuleEnable bool `json:"alarm_rule_enable"`

	// 告警状态。 - “OK”：正常 - “alarm”：超限阈值 - “Effective”：生效中 - “Invalid”：停用中
	AlarmRuleStatus string `json:"alarm_rule_status"`

	// 规则类型。 - “metric”：指标告警规则 - “event”：事件告警规则
	AlarmRuleType AlarmParamForV4DbAlarmRuleType `json:"alarm_rule_type"`

	MetricAlarmSpec *MetricAlarmSpec `json:"metric_alarm_spec,omitempty"`

	EventAlarmSpec *EventAlarmSpec `json:"event_alarm_spec,omitempty"`

	AlarmNotifications *AlarmNotification `json:"alarm_notifications"`

	// 用户id。
	UserId *string `json:"user_id,omitempty"`
}

func (o AlarmParamForV4Db) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AlarmParamForV4Db struct{}"
	}

	return strings.Join([]string{"AlarmParamForV4Db", string(data)}, " ")
}

type AlarmParamForV4DbAlarmRuleType struct {
	value string
}

type AlarmParamForV4DbAlarmRuleTypeEnum struct {
	METRIC AlarmParamForV4DbAlarmRuleType
	EVENT  AlarmParamForV4DbAlarmRuleType
}

func GetAlarmParamForV4DbAlarmRuleTypeEnum() AlarmParamForV4DbAlarmRuleTypeEnum {
	return AlarmParamForV4DbAlarmRuleTypeEnum{
		METRIC: AlarmParamForV4DbAlarmRuleType{
			value: "metric",
		},
		EVENT: AlarmParamForV4DbAlarmRuleType{
			value: "event",
		},
	}
}

func (c AlarmParamForV4DbAlarmRuleType) Value() string {
	return c.value
}

func (c AlarmParamForV4DbAlarmRuleType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AlarmParamForV4DbAlarmRuleType) UnmarshalJSON(b []byte) error {
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
