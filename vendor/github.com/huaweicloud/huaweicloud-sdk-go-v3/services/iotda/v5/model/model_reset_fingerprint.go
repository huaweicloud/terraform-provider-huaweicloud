package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResetFingerprint struct {

	// **参数说明**：设备指纹。设置改字段时平台将设备指纹重置为指定值；不携带时将之置空，后续设备第一次接入时，该设备指纹的值将设置为第一次接入时的证书指纹。 **取值范围**：长度为40的十六进制字符串或者长度为64的十六进制字符串。
	Fingerprint *string `json:"fingerprint,omitempty"`

	// **参数说明**：是否强制断开设备的连接，当前仅限长连接。默认值false。
	ForceDisconnect *bool `json:"force_disconnect,omitempty"`
}

func (o ResetFingerprint) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetFingerprint struct{}"
	}

	return strings.Join([]string{"ResetFingerprint", string(data)}, " ")
}
