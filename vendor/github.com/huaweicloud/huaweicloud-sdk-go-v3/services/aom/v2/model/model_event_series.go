package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 事件或者告警统计值统计结果元数据。
type EventSeries struct {

	// 事件或者告警级别枚举类型。
	EventSeverity *EventSeriesEventSeverity `json:"event_severity,omitempty"`

	// 事件或者告警统计结果。
	Values *[]int32 `json:"values,omitempty"`
}

func (o EventSeries) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventSeries struct{}"
	}

	return strings.Join([]string{"EventSeries", string(data)}, " ")
}

type EventSeriesEventSeverity struct {
	value string
}

type EventSeriesEventSeverityEnum struct {
	CRITICAL EventSeriesEventSeverity
	MAJOR    EventSeriesEventSeverity
	MINOR    EventSeriesEventSeverity
	INFO     EventSeriesEventSeverity
}

func GetEventSeriesEventSeverityEnum() EventSeriesEventSeverityEnum {
	return EventSeriesEventSeverityEnum{
		CRITICAL: EventSeriesEventSeverity{
			value: "Critical",
		},
		MAJOR: EventSeriesEventSeverity{
			value: "Major",
		},
		MINOR: EventSeriesEventSeverity{
			value: "Minor",
		},
		INFO: EventSeriesEventSeverity{
			value: "Info",
		},
	}
}

func (c EventSeriesEventSeverity) Value() string {
	return c.value
}

func (c EventSeriesEventSeverity) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EventSeriesEventSeverity) UnmarshalJSON(b []byte) error {
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
