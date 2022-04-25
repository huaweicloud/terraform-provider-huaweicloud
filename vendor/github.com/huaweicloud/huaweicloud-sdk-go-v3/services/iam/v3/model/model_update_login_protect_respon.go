package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// {  \"login_protect\":{         \"user_id\": \"16b26081f43d4c628c4bb88cf32e9f9b\",         \"enabled\": true,         \"verification_method\": \"vmfa\"     } }
type UpdateLoginProtectRespon struct {

	// 待修改信息的IAM用户ID。
	UserId string `json:"user_id"`

	// IAM用户是否开启登录保护，开启为\"true\"，不开启为\"false\"。
	Enabled bool `json:"enabled"`

	// IAM用户登录验证方式。手机验证为“sms”,邮箱验证为“email”,MFA验证为“vmfa”。
	VerificationMethod string `json:"verification_method"`
}

func (o UpdateLoginProtectRespon) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateLoginProtectRespon struct{}"
	}

	return strings.Join([]string{"UpdateLoginProtectRespon", string(data)}, " ")
}
