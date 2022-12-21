package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Notification struct {

	// 消息事件的名称.
	EventName *string `json:"event_name,omitempty"`

	// 事件通知模板选中状态
	Status *NotificationStatus `json:"status,omitempty"`

	// 事件通知主题的URN.
	Topic *string `json:"topic,omitempty"`

	// 订阅消息类型.
	MsgType *int32 `json:"msg_type,omitempty"`
}

func (o Notification) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Notification struct{}"
	}

	return strings.Join([]string{"Notification", string(data)}, " ")
}

type NotificationStatus struct {
	value string
}

type NotificationStatusEnum struct {
	ON  NotificationStatus
	OFF NotificationStatus
}

func GetNotificationStatusEnum() NotificationStatusEnum {
	return NotificationStatusEnum{
		ON: NotificationStatus{
			value: "on",
		},
		OFF: NotificationStatus{
			value: "off",
		},
	}
}

func (c NotificationStatus) Value() string {
	return c.value
}

func (c NotificationStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NotificationStatus) UnmarshalJSON(b []byte) error {
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
