package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type TokenUserResult struct {

	// IAM用户名。
	Name string `json:"name"`

	// IAM用户ID。
	Id string `json:"id"`

	// 密码过期时间（UTC时间），“”表示密码不过期。
	PasswordExpiresAt string `json:"password_expires_at"`

	Domain *TokenUserDomainResult `json:"domain"`
}

func (o TokenUserResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenUserResult struct{}"
	}

	return strings.Join([]string{"TokenUserResult", string(data)}, " ")
}
