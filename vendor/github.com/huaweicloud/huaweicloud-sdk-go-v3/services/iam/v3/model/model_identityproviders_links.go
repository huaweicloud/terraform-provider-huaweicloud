package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type IdentityprovidersLinks struct {

	// 身份提供商的资源链接地址。
	Self string `json:"self"`

	// 协议的资源链接地址。
	Protocols string `json:"protocols"`
}

func (o IdentityprovidersLinks) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IdentityprovidersLinks struct{}"
	}

	return strings.Join([]string{"IdentityprovidersLinks", string(data)}, " ")
}
