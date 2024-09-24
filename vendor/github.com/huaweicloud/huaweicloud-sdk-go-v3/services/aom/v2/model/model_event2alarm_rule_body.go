package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Event2alarmRuleBody 事件类告警实体
type Event2alarmRuleBody struct {

	// 用户项目id
	UserId string `json:"user_id"`

	// 规则名称。规则名称包含大小写字母，数字，特殊字符（_-）和汉字组成，不能以特殊字符开头或结尾，最大长度为100。
	Name string `json:"name"`

	// 规则描述。描述包含大小写字母，数字，特殊字符（_-<>=,.）和汉字组成，不能以下划线、中划线开头结尾，最大长度为1024。
	Description *string `json:"description,omitempty"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime *int64 `json:"update_time,omitempty"`

	// 事件源
	ResourceProvider *string `json:"resource_provider,omitempty"`

	Metadata *Event2alarmRuleBodyMetadata `json:"metadata"`

	// 规则是否启用
	Enable bool `json:"enable"`

	// 触发策略
	TriggerPolicies []Event2alarmRuleBodyTriggerPolicies `json:"trigger_policies"`

	// 告警类型。notification：直接告警。denoising：告警降噪。
	AlarmType Event2alarmRuleBodyAlarmType `json:"alarm_type"`

	// 告警行动规则
	ActionRule string `json:"action_rule"`

	// 告警抑制规则
	InhibitRule *string `json:"inhibit_rule,omitempty"`

	// 告警分组规则
	RouteGroupRule *string `json:"route_group_rule,omitempty"`

	// 事件名称
	EventNames *[]string `json:"event_names,omitempty"`

	// 是否迁移到2.0
	Migrated *bool `json:"migrated,omitempty"`

	// smn信息
	Topics *[]SmnTopics `json:"topics,omitempty"`
}

func (o Event2alarmRuleBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Event2alarmRuleBody struct{}"
	}

	return strings.Join([]string{"Event2alarmRuleBody", string(data)}, " ")
}

type Event2alarmRuleBodyAlarmType struct {
	value string
}

type Event2alarmRuleBodyAlarmTypeEnum struct {
	NOTIFICATION Event2alarmRuleBodyAlarmType
	DENOISING    Event2alarmRuleBodyAlarmType
}

func GetEvent2alarmRuleBodyAlarmTypeEnum() Event2alarmRuleBodyAlarmTypeEnum {
	return Event2alarmRuleBodyAlarmTypeEnum{
		NOTIFICATION: Event2alarmRuleBodyAlarmType{
			value: "notification",
		},
		DENOISING: Event2alarmRuleBodyAlarmType{
			value: "denoising",
		},
	}
}

func (c Event2alarmRuleBodyAlarmType) Value() string {
	return c.value
}

func (c Event2alarmRuleBodyAlarmType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *Event2alarmRuleBodyAlarmType) UnmarshalJSON(b []byte) error {
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
