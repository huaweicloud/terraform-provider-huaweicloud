package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResetFingerprint struct {

	// **参数说明**：设备指纹。设置该字段时平台将设备指纹重置为指定值；不携带时将之置空，后续设备第一次接入时，该设备指纹的值将设置为第一次接入时的证书指纹。 **取值范围**：长度为40的十六进制字符串或者长度为64的十六进制字符串。
	Fingerprint *string `json:"fingerprint,omitempty"`

	// **参数说明**：是否强制断开设备的连接，当前仅限长连接。默认值false。
	ForceDisconnect *bool `json:"force_disconnect,omitempty"`

	// **参数说明**：重置设备证书指纹的的类型。 **取值范围**： - PRIMARY：重置主指纹。设备证书鉴权优先使用的指纹，当设备接入物联网平台时，平台将优先使用主指纹进行校验。 - SECONDARY：重置辅指纹。设备的备用指纹，当主指纹校验不通过时，会启用辅指纹校验，辅指纹与主指纹有相同的效力。
	FingerprintType *string `json:"fingerprint_type,omitempty"`
}

func (o ResetFingerprint) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetFingerprint struct{}"
	}

	return strings.Join([]string{"ResetFingerprint", string(data)}, " ")
}
