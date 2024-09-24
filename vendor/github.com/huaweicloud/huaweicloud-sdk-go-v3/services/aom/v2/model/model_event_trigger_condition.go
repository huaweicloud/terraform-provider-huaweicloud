package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type EventTriggerCondition struct {

	// 事件名称。
	EventName *string `json:"event_name,omitempty"`

	// 触发方式： - “immediately”：立即触发 - “accumulative”：累计触发
	TriggerType *EventTriggerConditionTriggerType `json:"trigger_type,omitempty"`

	// 统计周期。单位为秒，例如 1小时 填“3600”，当trigger_type为“immediately”时 不填。
	AggregationWindow *int64 `json:"aggregation_window,omitempty"`

	// 判断条件：“>”,“<”,“=”,“>=”,“<=”，当trigger_type为“immediately”时 不填。
	Operator *string `json:"operator,omitempty"`

	// 键值对形式，键为告警级别，值为累计次数，当trigger_type为“immediately”时 值为空。
	Thresholds map[string]int32 `json:"thresholds,omitempty"`

	// 事件类告警频率。当trigger_type为“immediately”时 不填。 - “0”：只告警一次 - “300”：每5分钟 - “600”：每10分钟： - “900”：每15分钟： - “1800”：每30分钟： - “3600”：每1小时： - “10800”：每3小时： - “21600”：每6小时： - “43200”：每12小时： - “86400”：每天：
	Frequency *string `json:"frequency,omitempty"`
}

func (o EventTriggerCondition) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventTriggerCondition struct{}"
	}

	return strings.Join([]string{"EventTriggerCondition", string(data)}, " ")
}

type EventTriggerConditionTriggerType struct {
	value string
}

type EventTriggerConditionTriggerTypeEnum struct {
	IMMEDIATELY  EventTriggerConditionTriggerType
	ACCUMULATIVE EventTriggerConditionTriggerType
}

func GetEventTriggerConditionTriggerTypeEnum() EventTriggerConditionTriggerTypeEnum {
	return EventTriggerConditionTriggerTypeEnum{
		IMMEDIATELY: EventTriggerConditionTriggerType{
			value: "immediately",
		},
		ACCUMULATIVE: EventTriggerConditionTriggerType{
			value: "accumulative",
		},
	}
}

func (c EventTriggerConditionTriggerType) Value() string {
	return c.value
}

func (c EventTriggerConditionTriggerType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EventTriggerConditionTriggerType) UnmarshalJSON(b []byte) error {
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
