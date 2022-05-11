package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 阈值规则查询参数。
type QueryAlarmResult struct {

	// 是否启用通知。
	ActionEnabled *bool `json:"action_enabled,omitempty"`

	// 告警状态通知列表。
	AlarmActions *[]string `json:"alarm_actions,omitempty"`

	// 告警清除建议。
	AlarmAdvice *string `json:"alarm_advice,omitempty"`

	// 阈值规则描述。
	AlarmDescription *string `json:"alarm_description,omitempty"`

	// 告警级别。
	AlarmLevel *string `json:"alarm_level,omitempty"`

	// 阈值规则ID。
	AlarmRuleId *string `json:"alarm_rule_id,omitempty"`

	// 阈值规则名称。
	AlarmRuleName *string `json:"alarm_rule_name,omitempty"`

	// 极限条件。
	ComparisonOperator *string `json:"comparison_operator,omitempty"`

	// 时间序列维度。
	Dimensions *[]Dimension `json:"dimensions,omitempty"`

	// 间隔周期。
	EvaluationPeriods *int32 `json:"evaluation_periods,omitempty"`

	// 阈值规则是否启用。
	IdTurnOn *bool `json:"id_turn_on,omitempty"`

	// 数据不足通知列表。
	InsufficientDataActions *[]string `json:"insufficient_data_actions,omitempty"`

	// 时间序列名称。
	MetricName *string `json:"metric_name,omitempty"`

	// 时间序列命名空间。
	Namespace *string `json:"namespace,omitempty"`

	// 正常状态通知列表。
	OkActions *[]string `json:"ok_actions,omitempty"`

	// 统计周期。
	Period *int32 `json:"period,omitempty"`

	// 阈值规则模板名称。
	PolicyName *string `json:"policy_name,omitempty"`

	// 资源信息(已废弃)。
	Resources *[]string `json:"resources,omitempty"`

	// 原因描述。
	StateReason *string `json:"state_reason,omitempty"`

	// 状态更新时间戳。
	StateUpdatedTimestamp *string `json:"state_updated_timestamp,omitempty"`

	// 服务状态。
	StateValue *string `json:"state_value,omitempty"`

	// 统计方式。
	Statistic *string `json:"statistic,omitempty"`

	// 临界值。
	Threshold *string `json:"threshold,omitempty"`

	// 阈值规则类型。
	Type *string `json:"type,omitempty"`

	// 阈值单元。
	Unit *string `json:"unit,omitempty"`
}

func (o QueryAlarmResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueryAlarmResult struct{}"
	}

	return strings.Join([]string{"QueryAlarmResult", string(data)}, " ")
}
