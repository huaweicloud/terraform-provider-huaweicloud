package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UnbindMfaDevice struct {

	// 待解绑MFA设备的IAM用户ID。
	UserId string `json:"user_id"`

	// • 管理员为IAM用户解绑MFA设备：填写6位任意验证码，不做校验。 • IAM用户为自己解绑MFA设备：填写虚拟MFA验证码。
	AuthenticationCode string `json:"authentication_code"`

	// MFA设备序列号。
	SerialNumber string `json:"serial_number"`
}

func (o UnbindMfaDevice) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnbindMfaDevice struct{}"
	}

	return strings.Join([]string{"UnbindMfaDevice", string(data)}, " ")
}
