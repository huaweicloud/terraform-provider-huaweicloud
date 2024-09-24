package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListEventsRequest Request Object
type ListEventsRequest struct {

	// 查询类型。type=active_alert代表查询活动告警，type=history_alert代表查询历史告警。不传或者传其他值则返回指定查询条件的所有信息。
	Type *ListEventsRequestType `json:"type,omitempty"`

	// 企业项目id。 - 查询单个企业项目下实例，填写企业项目id。 - 查询所有企业项目下实例，填写“all_granted_eps”。
	EnterpriseProjectId *string `json:"Enterprise-Project-Id,omitempty"`

	// 不填默认值为1000
	Limit *int32 `json:"limit,omitempty"`

	// 分页标记，初始为0，后续值为返回体中的next_marker
	Marker *string `json:"marker,omitempty"`

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
