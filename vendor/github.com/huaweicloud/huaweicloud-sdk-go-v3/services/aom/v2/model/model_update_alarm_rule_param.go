package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 阈值规则实体
type UpdateAlarmRuleParam struct {

	// 是否启用通知。
	ActionEnabled *bool `json:"action_enabled,omitempty"`

	// 告警状态通知列表。
	AlarmActions *[]string `json:"alarm_actions,omitempty"`

	// 告警清除建议。
	AlarmAdvice *string `json:"alarm_advice,omitempty"`

	// 阈值规则描述。
	AlarmDescription *string `json:"alarm_description,omitempty"`

	// 告警级别。1：紧急，2：重要，3：一般，4：提示。
	AlarmLevel *UpdateAlarmRuleParamAlarmLevel `json:"alarm_level,omitempty"`

	// 阈值规则名称。
	AlarmRuleName string `json:"alarm_rule_name"`

	// 超限条件。
	ComparisonOperator *string `json:"comparison_operator,omitempty"`

	// 时间序列维度。
	Dimensions *[]Dimension `json:"dimensions,omitempty"`

	// 间隔周期。
	EvaluationPeriods *int32 `json:"evaluation_periods,omitempty"`

	// 阈值规则是否启用。
	IdTurnOn *bool `json:"id_turn_on,omitempty"`

	// 数据不足通知列表。
	InsufficientDataActions *[]string `json:"insufficient_data_actions,omitempty"`

	// 时间序列名称。名称长度取值范围为1~255个字符。
	MetricName *string `json:"metric_name,omitempty"`

	// 时间序列命名空间。
	Namespace *string `json:"namespace,omitempty"`

	// 正常状态通知列表。
	OkActions *[]string `json:"ok_actions,omitempty"`

	// 统计周期。
	Period *int32 `json:"period,omitempty"`

	// 统计方式。
	Statistic *UpdateAlarmRuleParamStatistic `json:"statistic,omitempty"`

	// 超限值。
	Threshold *string `json:"threshold,omitempty"`

	// 时间序列单位
	Unit *string `json:"unit,omitempty"`
}

func (o UpdateAlarmRuleParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAlarmRuleParam struct{}"
	}

	return strings.Join([]string{"UpdateAlarmRuleParam", string(data)}, " ")
}

type UpdateAlarmRuleParamAlarmLevel struct {
	value int32
}

type UpdateAlarmRuleParamAlarmLevelEnum struct {
	E_1 UpdateAlarmRuleParamAlarmLevel
	E_2 UpdateAlarmRuleParamAlarmLevel
	E_3 UpdateAlarmRuleParamAlarmLevel
	E_4 UpdateAlarmRuleParamAlarmLevel
}

func GetUpdateAlarmRuleParamAlarmLevelEnum() UpdateAlarmRuleParamAlarmLevelEnum {
	return UpdateAlarmRuleParamAlarmLevelEnum{
		E_1: UpdateAlarmRuleParamAlarmLevel{
			value: 1,
		}, E_2: UpdateAlarmRuleParamAlarmLevel{
			value: 2,
		}, E_3: UpdateAlarmRuleParamAlarmLevel{
			value: 3,
		}, E_4: UpdateAlarmRuleParamAlarmLevel{
			value: 4,
		},
	}
}

func (c UpdateAlarmRuleParamAlarmLevel) Value() int32 {
	return c.value
}

func (c UpdateAlarmRuleParamAlarmLevel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateAlarmRuleParamAlarmLevel) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("int32")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(int32)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to int32 error")
	}
}

type UpdateAlarmRuleParamStatistic struct {
	value string
}

type UpdateAlarmRuleParamStatisticEnum struct {
	MAXIMUM      UpdateAlarmRuleParamStatistic
	MINIMUM      UpdateAlarmRuleParamStatistic
	AVERAGE      UpdateAlarmRuleParamStatistic
	SUM          UpdateAlarmRuleParamStatistic
	SAMPLE_COUNT UpdateAlarmRuleParamStatistic
}

func GetUpdateAlarmRuleParamStatisticEnum() UpdateAlarmRuleParamStatisticEnum {
	return UpdateAlarmRuleParamStatisticEnum{
		MAXIMUM: UpdateAlarmRuleParamStatistic{
			value: "maximum",
		},
		MINIMUM: UpdateAlarmRuleParamStatistic{
			value: "minimum",
		},
		AVERAGE: UpdateAlarmRuleParamStatistic{
			value: "average",
		},
		SUM: UpdateAlarmRuleParamStatistic{
			value: "sum",
		},
		SAMPLE_COUNT: UpdateAlarmRuleParamStatistic{
			value: "sampleCount",
		},
	}
}

func (c UpdateAlarmRuleParamStatistic) Value() string {
	return c.value
}

func (c UpdateAlarmRuleParamStatistic) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateAlarmRuleParamStatistic) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
