package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type ProtocolLinks struct {

	// 身份提供商的资源链接地址。
	IdentityProvider string `json:"identity_provider"`

	// 资源链接地址。
	Self string `json:"self"`
}

func (o ProtocolLinks) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtocolLinks struct{}"
	}

	return strings.Join([]string{"ProtocolLinks", string(data)}, " ")
}
