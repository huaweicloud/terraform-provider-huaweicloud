package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateTokenWithIdTokenRequest struct {

	// 身份提供商ID。
	XIdpId string `json:"X-Idp-Id"`

	Body *GetIdTokenRequestBody `json:"body,omitempty"`
}

func (o CreateTokenWithIdTokenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTokenWithIdTokenRequest struct{}"
	}

	return strings.Join([]string{"CreateTokenWithIdTokenRequest", string(data)}, " ")
}
