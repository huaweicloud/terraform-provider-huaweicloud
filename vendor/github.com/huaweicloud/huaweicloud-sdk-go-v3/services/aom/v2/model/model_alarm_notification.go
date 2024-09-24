package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// AlarmNotification 告警规则通知模块。
type AlarmNotification struct {

	// 通知类型。 - “direct”：直接告警 - “alarm_policy”：告警降噪
	NotificationType AlarmNotificationNotificationType `json:"notification_type"`

	// 启用分组规则。 - 当通知类型为“alarm_policy”时：true - 当通知类型为“direct”时：false
	RouteGroupEnable bool `json:"route_group_enable"`

	// 分组规则名称。 - 当route_group_enable 为true时，填分组规则名称 - 当route_group_enable 为false时，填“”
	RouteGroupRule string `json:"route_group_rule"`

	// 是否启用告警行动规则。 - 当通知类型为“direct”时，填true - 当通知类型为“alarm_policy”时，填false
	NotificationEnable *bool `json:"notification_enable,omitempty"`

	// 告警行动策略id。 - 当notification_enable为true时，填告警行动策略id - 当notification_enable为false时，填“”
	BindNotificationRuleId *string `json:"bind_notification_rule_id,omitempty"`

	// 告警解决是否通知。 - true：通知 - false：不通知
	NotifyResolved *bool `json:"notify_resolved,omitempty"`

	// 告警触发是否通知。 - true：通知 - false：不通知
	NotifyTriggered *bool `json:"notify_triggered,omitempty"`

	// 通知频率 - 当通知类型为“alarm_policy”时，填“-1” - 当通知类型为“direct”时，    - “0”：只告警一次    - “300”：每5分钟    - “600”：每10分钟    - “900”：每15分钟    - “1800”：每30分钟    - “3600”：每1小时    - “10800”：每3小时    - “21600”：每6小时    - “43200”：每12小时    - “86400”：每天
	NotifyFrequency *int32 `json:"notify_frequency,omitempty"`
}

func (o AlarmNotification) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AlarmNotification struct{}"
	}

	return strings.Join([]string{"AlarmNotification", string(data)}, " ")
}

type AlarmNotificationNotificationType struct {
	value string
}

type AlarmNotificationNotificationTypeEnum struct {
	DIRECT       AlarmNotificationNotificationType
	ALARM_POLICY AlarmNotificationNotificationType
}

func GetAlarmNotificationNotificationTypeEnum() AlarmNotificationNotificationTypeEnum {
	return AlarmNotificationNotificationTypeEnum{
		DIRECT: AlarmNotificationNotificationType{
			value: "direct",
		},
		ALARM_POLICY: AlarmNotificationNotificationType{
			value: "alarm_policy",
		},
	}
}

func (c AlarmNotificationNotificationType) Value() string {
	return c.value
}

func (c AlarmNotificationNotificationType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AlarmNotificationNotificationType) UnmarshalJSON(b []byte) error {
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
