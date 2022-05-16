package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneUpdateIdentityProviderRequest struct {

	// 待更新的身份提供商ID。
	Id string `json:"id"`

	Body *KeystoneUpdateIdentityProviderRequestBody `json:"body,omitempty"`
}

func (o KeystoneUpdateIdentityProviderRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateIdentityProviderRequest struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateIdentityProviderRequest", string(data)}, " ")
}
