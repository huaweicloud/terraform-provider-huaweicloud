package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ResetFingerprintResponse struct {

	// 设备ID，用于唯一标识一个设备。在注册设备时直接指定，或者由物联网平台分配获得。由物联网平台分配时，生成规则为\"product_id\" + \"_\" + \"node_id\"拼接而成。
	DeviceId *string `json:"device_id,omitempty"`

	// 设备指纹。
	Fingerprint    *string `json:"fingerprint,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ResetFingerprintResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetFingerprintResponse struct{}"
	}

	return strings.Join([]string{"ResetFingerprintResponse", string(data)}, " ")
}
