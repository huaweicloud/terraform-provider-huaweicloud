package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 设备接入类型不返回密钥模式。
type AuthInfoWithoutSecret struct {

	// **参数说明**：指设备是否通过安全协议方式接入。 **取值范围**： - true：通过安全协议方式接入。 - false：通过非安全协议方式接入。非安全接入的设备存在被仿冒等安全风险，请谨慎使用。
	SecureAccess *bool `json:"secure_access,omitempty"`

	// **参数说明**：设备接入的有效时间，单位：秒，默认值：0。若设备在有效时间内未接入物联网平台并激活，则平台会删除该设备的注册信息。若设置为“0”，则表示平台不会删除该设备的注册信息（建议填写为“0”）。
	Timeout *int32 `json:"timeout,omitempty"`
}

func (o AuthInfoWithoutSecret) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AuthInfoWithoutSecret struct{}"
	}

	return strings.Join([]string{"AuthInfoWithoutSecret", string(data)}, " ")
}
