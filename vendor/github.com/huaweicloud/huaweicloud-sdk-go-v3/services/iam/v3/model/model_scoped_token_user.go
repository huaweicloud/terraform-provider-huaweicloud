package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type ScopedTokenUser struct {
	Domain *TokenDomainResult `json:"domain"`

	OsFederation *TokenUserOsfederation `json:"OS-FEDERATION"`

	// 用户ID。
	Id string `json:"id"`

	// 用户名。
	Name string `json:"name"`

	// 密码过期时间（UTC时间），“”表示密码不过期。
	PasswordExpiresAt string `json:"password_expires_at"`
}

func (o ScopedTokenUser) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScopedTokenUser struct{}"
	}

	return strings.Join([]string{"ScopedTokenUser", string(data)}, " ")
}
