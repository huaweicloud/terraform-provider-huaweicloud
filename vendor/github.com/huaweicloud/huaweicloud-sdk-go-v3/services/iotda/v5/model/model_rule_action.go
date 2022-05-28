package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则动作结构体
type RuleAction struct {

	// **参数说明**：规则动作的类型。 **取值范围**： - DEVICE_CMD：下发设备命令消息类型。 - SMN_FORWARDING：发送SMN消息类型。 - DEVICE_ALARM：上报设备告警消息类型。当选择该类型时，condition中必须有DEVICE_DATA条件类型。该类型动作只能唯一。 - DIS_FORWARDING：转发DIS服务消息类型。 - OBS_FORWARDING：转发OBS服务消息类型。 - ROMA_FORWARDING：转发ROMA Connect服务消息类型。 - IoTA_FORWARDING：转发IoTA服务消息类型。 - KAFKA_FORWARDING：转发kafka消息类型。
	Type string `json:"type"`

	SmnForwarding *ActionSmnForwarding `json:"smn_forwarding,omitempty"`

	DeviceAlarm *ActionDeviceAlarm `json:"device_alarm,omitempty"`

	DeviceCommand *ActionDeviceCommand `json:"device_command,omitempty"`
}

func (o RuleAction) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RuleAction struct{}"
	}

	return strings.Join([]string{"RuleAction", string(data)}, " ")
}
