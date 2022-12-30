package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则动作结构体
type RuleAction struct {

	// **参数说明**：规则动作的类型，端侧执行只支持下发设备命令消息类型。 **取值范围**： - DEVICE_CMD：下发设备命令消息类型。 - SMN_FORWARDING：发送SMN消息类型。 - DEVICE_ALARM：上报设备告警消息类型。当选择该类型时，condition中必须有DEVICE_DATA条件类型。该类型动作只能唯一。
	Type string `json:"type"`

	DeviceCommand *ActionDeviceCommand `json:"device_command,omitempty"`

	SmnForwarding *ActionSmnForwarding `json:"smn_forwarding,omitempty"`

	DeviceAlarm *ActionDeviceAlarm `json:"device_alarm,omitempty"`
}

func (o RuleAction) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RuleAction struct{}"
	}

	return strings.Join([]string{"RuleAction", string(data)}, " ")
}
