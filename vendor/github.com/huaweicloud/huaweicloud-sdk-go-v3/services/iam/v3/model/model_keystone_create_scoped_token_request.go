package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneCreateScopedTokenRequest struct {
	Body *KeystoneCreateScopedTokenRequestBody `json:"body,omitempty"`
}

func (o KeystoneCreateScopedTokenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateScopedTokenRequest struct{}"
	}

	return strings.Join([]string{"KeystoneCreateScopedTokenRequest", string(data)}, " ")
}
