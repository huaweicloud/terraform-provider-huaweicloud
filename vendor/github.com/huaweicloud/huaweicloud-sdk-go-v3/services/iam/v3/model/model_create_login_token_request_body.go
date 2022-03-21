package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CreateLoginTokenRequestBody struct {
	Auth *LoginTokenAuth `json:"auth"`
}

func (o CreateLoginTokenRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateLoginTokenRequestBody struct{}"
	}

	return strings.Join([]string{"CreateLoginTokenRequestBody", string(data)}, " ")
}
