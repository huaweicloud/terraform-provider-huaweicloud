package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneDeleteIdentityProviderRequest struct {

	// 待删除的身份提供商ID。
	Id string `json:"id"`
}

func (o KeystoneDeleteIdentityProviderRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneDeleteIdentityProviderRequest struct{}"
	}

	return strings.Join([]string{"KeystoneDeleteIdentityProviderRequest", string(data)}, " ")
}
