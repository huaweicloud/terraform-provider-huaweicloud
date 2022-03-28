package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type LoginTokenAuth struct {
	Securitytoken *LoginTokenSecurityToken `json:"securitytoken"`
}

func (o LoginTokenAuth) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LoginTokenAuth struct{}"
	}

	return strings.Join([]string{"LoginTokenAuth", string(data)}, " ")
}
