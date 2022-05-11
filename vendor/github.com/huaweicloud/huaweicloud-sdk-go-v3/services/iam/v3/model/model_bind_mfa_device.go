package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type BindMfaDevice struct {

	// 待绑定MFA设备的IAM用户ID。
	UserId string `json:"user_id"`

	// MFA设备序列号。
	SerialNumber string `json:"serial_number"`

	// 第一组验证码。
	AuthenticationCodeFirst string `json:"authentication_code_first"`

	// 第二组验证码。
	AuthenticationCodeSecond string `json:"authentication_code_second"`
}

func (o BindMfaDevice) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BindMfaDevice struct{}"
	}

	return strings.Join([]string{"BindMfaDevice", string(data)}, " ")
}
