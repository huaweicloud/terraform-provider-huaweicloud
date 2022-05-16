package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// token详细信息。
type ScopedTokenInfo struct {

	// 过期时间。
	ExpiresAt string `json:"expires_at"`

	// 获取token的方式，联邦用户默认为mapped。
	Methods []string `json:"methods"`

	// 生成时间。
	IssuedAt string `json:"issued_at"`

	User *FederationUserBody `json:"user"`

	Domain *DomainInfo `json:"domain,omitempty"`

	Project *ProjectInfo `json:"project,omitempty"`

	// roles信息。
	Roles []ScopedTokenInfoRoles `json:"roles"`

	// catalog信息
	Catalog []UnscopedTokenInfoCatalog `json:"catalog"`
}

func (o ScopedTokenInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScopedTokenInfo struct{}"
	}

	return strings.Join([]string{"ScopedTokenInfo", string(data)}, " ")
}
