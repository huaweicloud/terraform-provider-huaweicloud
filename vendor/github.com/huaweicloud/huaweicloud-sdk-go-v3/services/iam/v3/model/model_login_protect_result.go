package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type LoginProtectResult struct {

	// IAM用户是否开启登录保护，开启为\"true\"，未开启为\"false\"。
	Enabled bool `json:"enabled"`

	// IAM用户ID。
	UserId string `json:"user_id"`

	// IAM用户登录验证方式。
	VerificationMethod string `json:"verification_method"`
}

func (o LoginProtectResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LoginProtectResult struct{}"
	}

	return strings.Join([]string{"LoginProtectResult", string(data)}, " ")
}
