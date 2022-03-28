package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCreateScopedTokenResponse struct {
	Token *ScopeTokenResult `json:"token,omitempty"`

	XSubjectToken  *string `json:"X-Subject-Token,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o KeystoneCreateScopedTokenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateScopedTokenResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCreateScopedTokenResponse", string(data)}, " ")
}
