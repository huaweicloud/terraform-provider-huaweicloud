package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateUnscopedTokenWithIdTokenRequest struct {

	// 身份提供商id。
	IdpId string `json:"idp_id"`

	// 协议id。
	ProtocolId string `json:"protocol_id"`

	// OpenID Connect身份提供商的ID Token，格式为Bearer {ID Token}。
	Authorization string `json:"Authorization"`
}

func (o CreateUnscopedTokenWithIdTokenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateUnscopedTokenWithIdTokenRequest struct{}"
	}

	return strings.Join([]string{"CreateUnscopedTokenWithIdTokenRequest", string(data)}, " ")
}
