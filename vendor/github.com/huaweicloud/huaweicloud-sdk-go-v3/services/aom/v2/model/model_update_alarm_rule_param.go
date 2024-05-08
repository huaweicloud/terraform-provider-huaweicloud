package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UpdateAlarmRuleParam 阈值规则实体
type UpdateAlarmRuleParam struct {

	// 是否启用通知。
	ActionEnabled *bool `json:"action_enabled,omitempty"`

	// 告警状态通知列表。
	AlarmActions *[]string `json:"alarm_actions,omitempty"`

	// 告警清除建议。字符长度为0-255。
	AlarmAdvice *string `json:"alarm_advice,omitempty"`

	// 阈值规则描述。字符长度为0-1024。
	AlarmDescription *string `json:"alarm_description,omitempty"`

	// 告警级别。1：紧急，2：重要，3：一般，4：提示。
	AlarmLevel UpdateAlarmRuleParamAlarmLevel `json:"alarm_level"`

	// 阈值规则名称。规则名称包含大小写字母、数字、特殊字符（-_）和汉字组成，不能以特殊字符开头或结尾，最大长度为100。
	AlarmRuleName string `json:"alarm_rule_name"`

	// 超限条件。<：小于阈值。>：大于阈值。<=：小于等于阈值。>=：大于等于阈值。
	ComparisonOperator UpdateAlarmRuleParamComparisonOperator `json:"comparison_operator"`

	// 时间序列维度。
	Dimensions []Dimension `json:"dimensions"`

	// 间隔周期。
	EvaluationPeriods int32 `json:"evaluation_periods"`

	// 阈值规则是否启用。
	IsTurnOn *bool `json:"is_turn_on,omitempty"`

	// 数据不足通知列表。
	InsufficientDataActions *[]string `json:"insufficient_data_actions,omitempty"`

	// 时间序列名称。名称长度取值范围为1~255个字符。
	MetricName string `json:"metric_name"`

	// 时间序列命名空间。
	Namespace string `json:"namespace"`

	// 正常状态通知列表。
	OkActions *[]string `json:"ok_actions,omitempty"`

	// 统计周期。60000：一分钟。300000：五分钟。900000：十五分钟。3600000：一小时。
	Period UpdateAlarmRuleParamPeriod `json:"period"`

	// 统计方式。
	Statistic UpdateAlarmRuleParamStatistic `json:"statistic"`

	// 超限值。
	Threshold string `json:"threshold"`

	// 时间序列单位
	Unit string `json:"unit"`
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
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: int32")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(int32); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to int32 error")
	}
}

type UpdateAlarmRuleParamComparisonOperator struct {
	value string
}

type UpdateAlarmRuleParamComparisonOperatorEnum struct {
	LESS_THAN                UpdateAlarmRuleParamComparisonOperator
	GREATER_THAN             UpdateAlarmRuleParamComparisonOperator
	LESS_THAN_OR_EQUAL_TO    UpdateAlarmRuleParamComparisonOperator
	GREATER_THAN_OR_EQUAL_TO UpdateAlarmRuleParamComparisonOperator
}

func GetUpdateAlarmRuleParamComparisonOperatorEnum() UpdateAlarmRuleParamComparisonOperatorEnum {
	return UpdateAlarmRuleParamComparisonOperatorEnum{
		LESS_THAN: UpdateAlarmRuleParamComparisonOperator{
			value: "<",
		},
		GREATER_THAN: UpdateAlarmRuleParamComparisonOperator{
			value: ">",
		},
		LESS_THAN_OR_EQUAL_TO: UpdateAlarmRuleParamComparisonOperator{
			value: "<=",
		},
		GREATER_THAN_OR_EQUAL_TO: UpdateAlarmRuleParamComparisonOperator{
			value: ">=",
		},
	}
}

func (c UpdateAlarmRuleParamComparisonOperator) Value() string {
	return c.value
}

func (c UpdateAlarmRuleParamComparisonOperator) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateAlarmRuleParamComparisonOperator) UnmarshalJSON(b []byte) error {
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

type UpdateAlarmRuleParamPeriod struct {
	value int32
}

type UpdateAlarmRuleParamPeriodEnum struct {
	E_60000   UpdateAlarmRuleParamPeriod
	E_300000  UpdateAlarmRuleParamPeriod
	E_900000  UpdateAlarmRuleParamPeriod
	E_3600000 UpdateAlarmRuleParamPeriod
}

func GetUpdateAlarmRuleParamPeriodEnum() UpdateAlarmRuleParamPeriodEnum {
	return UpdateAlarmRuleParamPeriodEnum{
		E_60000: UpdateAlarmRuleParamPeriod{
			value: 60000,
		}, E_300000: UpdateAlarmRuleParamPeriod{
			value: 300000,
		}, E_900000: UpdateAlarmRuleParamPeriod{
			value: 900000,
		}, E_3600000: UpdateAlarmRuleParamPeriod{
			value: 3600000,
		},
	}
}

func (c UpdateAlarmRuleParamPeriod) Value() int32 {
	return c.value
}

func (c UpdateAlarmRuleParamPeriod) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateAlarmRuleParamPeriod) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("int32")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: int32")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(int32); ok {
		c.value = val
		return nil
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
