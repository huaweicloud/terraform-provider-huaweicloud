package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 修改设备信息对象。
type UpdateDevice struct {

	// **参数说明**：设备名称。 **取值范围**：长度不超过256，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合，建议不少于4个字符。
	DeviceName *string `json:"device_name,omitempty"`

	// **参数说明**：设备的描述信息。 **取值范围**：长度不超过2048，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合
	Description *string `json:"description,omitempty"`

	// **参数说明**：设备扩展信息。用户可以自定义任何想要的扩展信息，修改子设备信息时不会下发给网关。
	ExtensionInfo *interface{} `json:"extension_info,omitempty"`

	AuthInfo *AuthInfoWithoutSecret `json:"auth_info,omitempty"`
}

func (o UpdateDevice) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDevice struct{}"
	}

	return strings.Join([]string{"UpdateDevice", string(data)}, " ")
}
