package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneCreateScopedTokenRequestBody struct {
	Auth *ScopedTokenAuth `json:"auth"`
}

func (o KeystoneCreateScopedTokenRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateScopedTokenRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneCreateScopedTokenRequestBody", string(data)}, " ")
}
