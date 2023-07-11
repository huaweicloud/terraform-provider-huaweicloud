package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Notifications 发送结果
type Notifications struct {

	// 告警行动规则名称
	ActionRule *string `json:"action_rule,omitempty"`

	// 通知类型。SMN：消息通知服务
	NotifierChannel *NotificationsNotifierChannel `json:"notifier_channel,omitempty"`

	SmnChannel *SmnResponse `json:"smn_channel,omitempty"`
}

func (o Notifications) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Notifications struct{}"
	}

	return strings.Join([]string{"Notifications", string(data)}, " ")
}

type NotificationsNotifierChannel struct {
	value string
}

type NotificationsNotifierChannelEnum struct {
	SMN NotificationsNotifierChannel
}

func GetNotificationsNotifierChannelEnum() NotificationsNotifierChannelEnum {
	return NotificationsNotifierChannelEnum{
		SMN: NotificationsNotifierChannel{
			value: "SMN",
		},
	}
}

func (c NotificationsNotifierChannel) Value() string {
	return c.value
}

func (c NotificationsNotifierChannel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NotificationsNotifierChannel) UnmarshalJSON(b []byte) error {
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
