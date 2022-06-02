package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListEventsRequest struct {

	// 查询类型。type=active_alert代表查询活动告警，type=history_alert代表查询历史告警。不传或者传其他值则返回指定查询条件的所有信息。
	Type *ListEventsRequestType `json:"type,omitempty"`

	Body *EventQueryParam2 `json:"body,omitempty"`
}

func (o ListEventsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEventsRequest struct{}"
	}

	return strings.Join([]string{"ListEventsRequest", string(data)}, " ")
}

type ListEventsRequestType struct {
	value string
}

type ListEventsRequestTypeEnum struct {
	HISTORY_ALERT ListEventsRequestType
	ACTIVE_ALERT  ListEventsRequestType
}

func GetListEventsRequestTypeEnum() ListEventsRequestTypeEnum {
	return ListEventsRequestTypeEnum{
		HISTORY_ALERT: ListEventsRequestType{
			value: "history_alert",
		},
		ACTIVE_ALERT: ListEventsRequestType{
			value: "active_alert",
		},
	}
}

func (c ListEventsRequestType) Value() string {
	return c.value
}

func (c ListEventsRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListEventsRequestType) UnmarshalJSON(b []byte) error {
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
