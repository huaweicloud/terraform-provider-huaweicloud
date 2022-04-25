package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// token详细信息。
type UnscopedTokenInfo struct {

	// 过期时间。
	ExpiresAt string `json:"expires_at"`

	// token获取方式，联邦认证默认为mapped。
	Methods []string `json:"methods"`

	// 生成时间。
	IssuedAt string `json:"issued_at"`

	User *FederationUserBody `json:"user"`

	// roles信息。
	Roles *[]UnscopedTokenInfoRoles `json:"roles,omitempty"`

	// catalog信息。
	Catalog *[]UnscopedTokenInfoCatalog `json:"catalog,omitempty"`
}

func (o UnscopedTokenInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnscopedTokenInfo struct{}"
	}

	return strings.Join([]string{"UnscopedTokenInfo", string(data)}, " ")
}
