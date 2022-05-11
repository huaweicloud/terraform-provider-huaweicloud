package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneListIdentityProvidersResponse struct {

	// 身份提供商信息列表。
	IdentityProviders *[]IdentityprovidersResult `json:"identity_providers,omitempty"`

	Links          *Links `json:"links,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o KeystoneListIdentityProvidersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListIdentityProvidersResponse struct{}"
	}

	return strings.Join([]string{"KeystoneListIdentityProvidersResponse", string(data)}, " ")
}
