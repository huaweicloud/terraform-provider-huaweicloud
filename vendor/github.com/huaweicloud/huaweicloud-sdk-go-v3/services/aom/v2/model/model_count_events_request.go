package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type CountEventsRequest struct {

	// 查询类型。type=active_alert代表查询活动告警，type=history_alert代表查询历史告警。不传或者传其他值则返回指定查询条件的所有信息。
	Type *CountEventsRequestType `json:"type,omitempty"`

	Body *EventQueryParam `json:"body,omitempty"`
}

func (o CountEventsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CountEventsRequest struct{}"
	}

	return strings.Join([]string{"CountEventsRequest", string(data)}, " ")
}

type CountEventsRequestType struct {
	value string
}

type CountEventsRequestTypeEnum struct {
	HISTORY_ALERT CountEventsRequestType
	ACTIVE_ALERT  CountEventsRequestType
}

func GetCountEventsRequestTypeEnum() CountEventsRequestTypeEnum {
	return CountEventsRequestTypeEnum{
		HISTORY_ALERT: CountEventsRequestType{
			value: "history_alert",
		},
		ACTIVE_ALERT: CountEventsRequestType{
			value: "active_alert",
		},
	}
}

func (c CountEventsRequestType) Value() string {
	return c.value
}

func (c CountEventsRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CountEventsRequestType) UnmarshalJSON(b []byte) error {
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
