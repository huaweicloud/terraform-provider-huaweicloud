package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResetDeviceSecret struct {

	// **参数说明**：设备密钥，设置该字段时平台将设备密钥重置为指定值，若不设置则由平台自动生成。 **取值范围**：长度不低于8不超过32，只允许字母、数字、下划线（_）、连接符（-）的组合。
	Secret *string `json:"secret,omitempty"`

	// **参数说明**：是否强制断开设备的连接，当前仅限长连接。默认值false。
	ForceDisconnect *bool `json:"force_disconnect,omitempty"`
}

func (o ResetDeviceSecret) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetDeviceSecret struct{}"
	}

	return strings.Join([]string{"ResetDeviceSecret", string(data)}, " ")
}
