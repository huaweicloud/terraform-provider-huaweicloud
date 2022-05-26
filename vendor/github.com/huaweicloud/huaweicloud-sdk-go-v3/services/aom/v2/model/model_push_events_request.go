package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type PushEventsRequest struct {

	// 告警所属的企业项目id。
	XEnterprisePrjectId *string `json:"x-enterprise-prject-id,omitempty"`

	// 接口请求动作。action=clear代表清除告警，不传或者传其他值默认为上报动作。
	Action *PushEventsRequestAction `json:"action,omitempty"`

	Body *EventList `json:"body,omitempty"`
}

func (o PushEventsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PushEventsRequest struct{}"
	}

	return strings.Join([]string{"PushEventsRequest", string(data)}, " ")
}

type PushEventsRequestAction struct {
	value string
}

type PushEventsRequestActionEnum struct {
	CLEAR PushEventsRequestAction
}

func GetPushEventsRequestActionEnum() PushEventsRequestActionEnum {
	return PushEventsRequestActionEnum{
		CLEAR: PushEventsRequestAction{
			value: "clear",
		},
	}
}

func (c PushEventsRequestAction) Value() string {
	return c.value
}

func (c PushEventsRequestAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *PushEventsRequestAction) UnmarshalJSON(b []byte) error {
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
