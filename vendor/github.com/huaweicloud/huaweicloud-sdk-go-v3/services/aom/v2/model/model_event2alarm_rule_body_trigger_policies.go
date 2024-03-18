package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Event2alarmRuleBodyTriggerPolicies struct {

	// 自增编号
	Id *int32 `json:"id,omitempty"`

	// 事件名称
	Name *string `json:"name,omitempty"`

	// 触发类型。accumulative: 累计触发，immediately: 立即触发
	TriggerType *Event2alarmRuleBodyTriggerPoliciesTriggerType `json:"trigger_type,omitempty"`

	// 触发周期
	Period *int32 `json:"period,omitempty"`

	// 比较符
	Operator *string `json:"operator,omitempty"`

	// 触发次数
	Count *int32 `json:"count,omitempty"`

	// 告警等级
	Level *string `json:"level,omitempty"`
}

func (o Event2alarmRuleBodyTriggerPolicies) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Event2alarmRuleBodyTriggerPolicies struct{}"
	}

	return strings.Join([]string{"Event2alarmRuleBodyTriggerPolicies", string(data)}, " ")
}

type Event2alarmRuleBodyTriggerPoliciesTriggerType struct {
	value string
}

type Event2alarmRuleBodyTriggerPoliciesTriggerTypeEnum struct {
	ACCUMULATIVE Event2alarmRuleBodyTriggerPoliciesTriggerType
	IMMEDIATELY  Event2alarmRuleBodyTriggerPoliciesTriggerType
}

func GetEvent2alarmRuleBodyTriggerPoliciesTriggerTypeEnum() Event2alarmRuleBodyTriggerPoliciesTriggerTypeEnum {
	return Event2alarmRuleBodyTriggerPoliciesTriggerTypeEnum{
		ACCUMULATIVE: Event2alarmRuleBodyTriggerPoliciesTriggerType{
			value: "accumulative",
		},
		IMMEDIATELY: Event2alarmRuleBodyTriggerPoliciesTriggerType{
			value: "immediately",
		},
	}
}

func (c Event2alarmRuleBodyTriggerPoliciesTriggerType) Value() string {
	return c.value
}

func (c Event2alarmRuleBodyTriggerPoliciesTriggerType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *Event2alarmRuleBodyTriggerPoliciesTriggerType) UnmarshalJSON(b []byte) error {
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
