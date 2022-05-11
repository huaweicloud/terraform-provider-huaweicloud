package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type ScopeTokenResult struct {

	// 获取token的方式。
	Methods []string `json:"methods"`

	// token过期时间。
	ExpiresAt string `json:"expires_at"`

	// 服务目录信息。
	Catalog *[]TokenCatalog `json:"catalog,omitempty"`

	Domain *TokenDomainResult `json:"domain"`

	Project *TokenProjectResult `json:"project"`

	// token的权限信息。
	Roles []TokenRole `json:"roles"`

	User *ScopedTokenUser `json:"user"`

	// token下发时间。
	IssuedAt string `json:"issued_at"`
}

func (o ScopeTokenResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScopeTokenResult struct{}"
	}

	return strings.Join([]string{"ScopeTokenResult", string(data)}, " ")
}
