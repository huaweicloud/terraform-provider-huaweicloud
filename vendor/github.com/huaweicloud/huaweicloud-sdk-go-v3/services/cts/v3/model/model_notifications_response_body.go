package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type NotificationsResponseBody struct {

	// 标识关键操作名称。
	NotificationName *string `json:"notification_name,omitempty"`

	// 标识操作类型。 目前支持的操作类型有完整类型(complete)和自定义类型(customized)。 完整类型下，CTS发送通知的对象为已对接服务的所有事件。 自定义类型下，CTS发送通知的对象是在operations列表中指定的事件。
	OperationType *NotificationsResponseBodyOperationType `json:"operation_type,omitempty"`

	// 云服务委托名称。
	AgencyName *NotificationsResponseBodyAgencyName `json:"agency_name,omitempty"`

	// 操作事件列表。
	Operations *[]Operations `json:"operations,omitempty"`

	// 通知用户列表，目前最多支持对10个用户组和50个用户发起的操作进行配置。
	NotifyUserList *[]NotificationUsers `json:"notify_user_list,omitempty"`

	// 标识关键操作通知状态，包括正常(enabled)，停止(disabled)两种状态。
	Status *NotificationsResponseBodyStatus `json:"status,omitempty"`

	// 消息通知服务的topic_urn或者函数工作流的func_urn。 - 消息通知服务的topic_urn可以通过消息通知服务的查询主题列表API获取，示例：urn:smn:regionId:f96188c7ccaf4ffba0c9aa149ab2bd57:test_topic_v2。 - 函数工作流的func_urn可以通过函数工作流的获取函数列表API获取，示例：urn:fss:xxxxxxxxx:7aad83af3e8d42e99ac194e8419e2c9b:function:default:test。
	TopicId *string `json:"topic_id,omitempty"`

	// 关键操作通知的唯一标识。
	NotificationId *string `json:"notification_id,omitempty"`

	// 关键操作通知类型，根据topic_id区分为消息通知服务(smn)和函数工作流(fun)。
	NotificationType *NotificationsResponseBodyNotificationType `json:"notification_type,omitempty"`

	// 项目ID。
	ProjectId *string `json:"project_id,omitempty"`

	// 关键操作通知创建时间戳。
	CreateTime *int64 `json:"create_time,omitempty"`

	Filter *Filter `json:"filter,omitempty"`
}

func (o NotificationsResponseBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NotificationsResponseBody struct{}"
	}

	return strings.Join([]string{"NotificationsResponseBody", string(data)}, " ")
}

type NotificationsResponseBodyOperationType struct {
	value string
}

type NotificationsResponseBodyOperationTypeEnum struct {
	CUSTOMIZED NotificationsResponseBodyOperationType
	COMPLETE   NotificationsResponseBodyOperationType
}

func GetNotificationsResponseBodyOperationTypeEnum() NotificationsResponseBodyOperationTypeEnum {
	return NotificationsResponseBodyOperationTypeEnum{
		CUSTOMIZED: NotificationsResponseBodyOperationType{
			value: "customized",
		},
		COMPLETE: NotificationsResponseBodyOperationType{
			value: "complete",
		},
	}
}

func (c NotificationsResponseBodyOperationType) Value() string {
	return c.value
}

func (c NotificationsResponseBodyOperationType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NotificationsResponseBodyOperationType) UnmarshalJSON(b []byte) error {
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

type NotificationsResponseBodyAgencyName struct {
	value string
}

type NotificationsResponseBodyAgencyNameEnum struct {
	CTS_ADMIN_TRUST NotificationsResponseBodyAgencyName
}

func GetNotificationsResponseBodyAgencyNameEnum() NotificationsResponseBodyAgencyNameEnum {
	return NotificationsResponseBodyAgencyNameEnum{
		CTS_ADMIN_TRUST: NotificationsResponseBodyAgencyName{
			value: "cts_admin_trust",
		},
	}
}

func (c NotificationsResponseBodyAgencyName) Value() string {
	return c.value
}

func (c NotificationsResponseBodyAgencyName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NotificationsResponseBodyAgencyName) UnmarshalJSON(b []byte) error {
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

type NotificationsResponseBodyStatus struct {
	value string
}

type NotificationsResponseBodyStatusEnum struct {
	ENABLED  NotificationsResponseBodyStatus
	DISABLED NotificationsResponseBodyStatus
}

func GetNotificationsResponseBodyStatusEnum() NotificationsResponseBodyStatusEnum {
	return NotificationsResponseBodyStatusEnum{
		ENABLED: NotificationsResponseBodyStatus{
			value: "enabled",
		},
		DISABLED: NotificationsResponseBodyStatus{
			value: "disabled",
		},
	}
}

func (c NotificationsResponseBodyStatus) Value() string {
	return c.value
}

func (c NotificationsResponseBodyStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NotificationsResponseBodyStatus) UnmarshalJSON(b []byte) error {
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

type NotificationsResponseBodyNotificationType struct {
	value string
}

type NotificationsResponseBodyNotificationTypeEnum struct {
	SMN NotificationsResponseBodyNotificationType
	FUN NotificationsResponseBodyNotificationType
}

func GetNotificationsResponseBodyNotificationTypeEnum() NotificationsResponseBodyNotificationTypeEnum {
	return NotificationsResponseBodyNotificationTypeEnum{
		SMN: NotificationsResponseBodyNotificationType{
			value: "smn",
		},
		FUN: NotificationsResponseBodyNotificationType{
			value: "fun",
		},
	}
}

func (c NotificationsResponseBodyNotificationType) Value() string {
	return c.value
}

func (c NotificationsResponseBodyNotificationType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NotificationsResponseBodyNotificationType) UnmarshalJSON(b []byte) error {
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
