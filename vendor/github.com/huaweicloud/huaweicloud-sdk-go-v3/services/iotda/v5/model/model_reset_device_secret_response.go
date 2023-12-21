package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResetDeviceSecretResponse Response Object
type ResetDeviceSecretResponse struct {

	// 设备ID，用于唯一标识一个设备。在注册设备时直接指定，或者由物联网平台分配获得。由物联网平台分配时，生成规则为\"product_id\" + \"_\" + \"node_id\"拼接而成。
	DeviceId *string `json:"device_id,omitempty"`

	// 设备密钥。
	Secret *string `json:"secret,omitempty"`

	// **参数说明**：重置设备秘钥的的类型。 **取值范围**： - PRIMARY：重置主秘钥。设备秘钥鉴权优先使用的密钥，当设备接入物联网平台时，平台将优先使用主密钥进行校验。 - SECONDARY：重置辅秘钥。设备的备用密钥，当主密钥校验不通过时，会启用辅密钥校验，辅密钥与主密钥有相同的效力；辅密钥对coap协议接入的设备不生效。
	SecretType     *string `json:"secret_type,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ResetDeviceSecretResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetDeviceSecretResponse struct{}"
	}

	return strings.Join([]string{"ResetDeviceSecretResponse", string(data)}, " ")
}
