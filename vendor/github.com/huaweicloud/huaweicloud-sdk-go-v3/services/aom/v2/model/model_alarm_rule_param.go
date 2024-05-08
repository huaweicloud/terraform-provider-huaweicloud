package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// AlarmRuleParam 阈值规则实体
type AlarmRuleParam struct {

	// 是否启用通知。
	ActionEnabled *bool `json:"action_enabled,omitempty"`

	// 告警状态通知列表。
	AlarmActions *[]string `json:"alarm_actions,omitempty"`

	// 告警清除建议。
	AlarmAdvice *string `json:"alarm_advice,omitempty"`

	// 阈值规则描述。
	AlarmDescription *string `json:"alarm_description,omitempty"`

	// 告警级别。1：紧急，2：重要，3：一般，4：提示。
	AlarmLevel AlarmRuleParamAlarmLevel `json:"alarm_level"`

	// 阈值规则名称。规则名称包含大小写字母、数字、特殊字符（-_）和汉字组成，不能以特殊字符开头或结尾，最大长度为100。
	AlarmRuleName string `json:"alarm_rule_name"`

	// 超限条件。<：小于阈值。>：大于阈值。<=：小于等于阈值。>=：大于等于阈值。
	ComparisonOperator AlarmRuleParamComparisonOperator `json:"comparison_operator"`

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
	Period AlarmRuleParamPeriod `json:"period"`

	// 统计方式。
	Statistic AlarmRuleParamStatistic `json:"statistic"`

	// 超限值。
	Threshold string `json:"threshold"`

	// 时间序列单位
	Unit string `json:"unit"`
}

func (o AlarmRuleParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AlarmRuleParam struct{}"
	}

	return strings.Join([]string{"AlarmRuleParam", string(data)}, " ")
}

type AlarmRuleParamAlarmLevel struct {
	value int32
}

type AlarmRuleParamAlarmLevelEnum struct {
	E_1 AlarmRuleParamAlarmLevel
	E_2 AlarmRuleParamAlarmLevel
	E_3 AlarmRuleParamAlarmLevel
	E_4 AlarmRuleParamAlarmLevel
}

func GetAlarmRuleParamAlarmLevelEnum() AlarmRuleParamAlarmLevelEnum {
	return AlarmRuleParamAlarmLevelEnum{
		E_1: AlarmRuleParamAlarmLevel{
			value: 1,
		}, E_2: AlarmRuleParamAlarmLevel{
			value: 2,
		}, E_3: AlarmRuleParamAlarmLevel{
			value: 3,
		}, E_4: AlarmRuleParamAlarmLevel{
			value: 4,
		},
	}
}

func (c AlarmRuleParamAlarmLevel) Value() int32 {
	return c.value
}

func (c AlarmRuleParamAlarmLevel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AlarmRuleParamAlarmLevel) UnmarshalJSON(b []byte) error {
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

type AlarmRuleParamComparisonOperator struct {
	value string
}

type AlarmRuleParamComparisonOperatorEnum struct {
	LESS_THAN                AlarmRuleParamComparisonOperator
	GREATER_THAN             AlarmRuleParamComparisonOperator
	LESS_THAN_OR_EQUAL_TO    AlarmRuleParamComparisonOperator
	GREATER_THAN_OR_EQUAL_TO AlarmRuleParamComparisonOperator
}

func GetAlarmRuleParamComparisonOperatorEnum() AlarmRuleParamComparisonOperatorEnum {
	return AlarmRuleParamComparisonOperatorEnum{
		LESS_THAN: AlarmRuleParamComparisonOperator{
			value: "<",
		},
		GREATER_THAN: AlarmRuleParamComparisonOperator{
			value: ">",
		},
		LESS_THAN_OR_EQUAL_TO: AlarmRuleParamComparisonOperator{
			value: "<=",
		},
		GREATER_THAN_OR_EQUAL_TO: AlarmRuleParamComparisonOperator{
			value: ">=",
		},
	}
}

func (c AlarmRuleParamComparisonOperator) Value() string {
	return c.value
}

func (c AlarmRuleParamComparisonOperator) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AlarmRuleParamComparisonOperator) UnmarshalJSON(b []byte) error {
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

type AlarmRuleParamPeriod struct {
	value int32
}

type AlarmRuleParamPeriodEnum struct {
	E_60000   AlarmRuleParamPeriod
	E_300000  AlarmRuleParamPeriod
	E_900000  AlarmRuleParamPeriod
	E_3600000 AlarmRuleParamPeriod
}

func GetAlarmRuleParamPeriodEnum() AlarmRuleParamPeriodEnum {
	return AlarmRuleParamPeriodEnum{
		E_60000: AlarmRuleParamPeriod{
			value: 60000,
		}, E_300000: AlarmRuleParamPeriod{
			value: 300000,
		}, E_900000: AlarmRuleParamPeriod{
			value: 900000,
		}, E_3600000: AlarmRuleParamPeriod{
			value: 3600000,
		},
	}
}

func (c AlarmRuleParamPeriod) Value() int32 {
	return c.value
}

func (c AlarmRuleParamPeriod) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AlarmRuleParamPeriod) UnmarshalJSON(b []byte) error {
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

type AlarmRuleParamStatistic struct {
	value string
}

type AlarmRuleParamStatisticEnum struct {
	MAXIMUM      AlarmRuleParamStatistic
	MINIMUM      AlarmRuleParamStatistic
	AVERAGE      AlarmRuleParamStatistic
	SUM          AlarmRuleParamStatistic
	SAMPLE_COUNT AlarmRuleParamStatistic
}

func GetAlarmRuleParamStatisticEnum() AlarmRuleParamStatisticEnum {
	return AlarmRuleParamStatisticEnum{
		MAXIMUM: AlarmRuleParamStatistic{
			value: "maximum",
		},
		MINIMUM: AlarmRuleParamStatistic{
			value: "minimum",
		},
		AVERAGE: AlarmRuleParamStatistic{
			value: "average",
		},
		SUM: AlarmRuleParamStatistic{
			value: "sum",
		},
		SAMPLE_COUNT: AlarmRuleParamStatistic{
			value: "sampleCount",
		},
	}
}

func (c AlarmRuleParamStatistic) Value() string {
	return c.value
}

func (c AlarmRuleParamStatistic) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AlarmRuleParamStatistic) UnmarshalJSON(b []byte) error {
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
