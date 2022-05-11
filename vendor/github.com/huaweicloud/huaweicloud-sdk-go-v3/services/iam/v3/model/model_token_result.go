package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type TokenResult struct {

	// 服务目录信息。
	Catalog []TokenCatalog `json:"catalog"`

	Domain *TokenDomainResult `json:"domain,omitempty"`

	// token过期时间。
	ExpiresAt string `json:"expires_at"`

	// token下发时间。
	IssuedAt string `json:"issued_at"`

	// 获取token的方式。
	Methods []string `json:"methods"`

	Project *TokenProjectResult `json:"project,omitempty"`

	// token的权限信息。
	Roles []TokenRole `json:"roles"`

	User *TokenUserResult `json:"user"`
}

func (o TokenResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenResult struct{}"
	}

	return strings.Join([]string{"TokenResult", string(data)}, " ")
}
