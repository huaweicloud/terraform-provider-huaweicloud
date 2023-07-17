package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Event2alarmRuleBody 事件类告警实体
type Event2alarmRuleBody struct {

	// 用户项目id
	UserId string `json:"user_id"`

	// 规则名称
	Name string `json:"name"`

	// 规则描述
	Description *string `json:"description,omitempty"`

	// 创建时间
	CreateTime int32 `json:"create_time"`

	// 更新时间
	UpdateTime *int32 `json:"update_time,omitempty"`

	// 事件源
	ResourceProvider *string `json:"resource_provider,omitempty"`

	Metadata *Event2alarmRuleBodyMetadata `json:"metadata"`

	// 规则是否启用
	Enable bool `json:"enable"`

	// 触发策略
	TriggerPolicies []Event2alarmRuleBodyTriggerPolicies `json:"trigger_policies"`

	// 告警类型
	AlarmType string `json:"alarm_type"`

	// 告警行动规则
	ActionRule string `json:"action_rule"`

	// 告警抑制规则
	InhibitRule *string `json:"inhibit_rule,omitempty"`

	// 告警静默规则
	RouteGroupRule *string `json:"route_group_rule,omitempty"`
}

func (o Event2alarmRuleBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Event2alarmRuleBody struct{}"
	}

	return strings.Join([]string{"Event2alarmRuleBody", string(data)}, " ")
}
