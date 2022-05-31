package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 下发设备命令消息结构
type ActionDeviceCommand struct {

	// **参数说明**：下发命令的设备ID。当创建设备数据规则时，若device_id为空，则命令下发给触发条件的设备。当创建定时规则时，不允许为空。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	DeviceId *string `json:"device_id,omitempty"`

	Cmd *Cmd `json:"cmd"`
}

func (o ActionDeviceCommand) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ActionDeviceCommand struct{}"
	}

	return strings.Join([]string{"ActionDeviceCommand", string(data)}, " ")
}
