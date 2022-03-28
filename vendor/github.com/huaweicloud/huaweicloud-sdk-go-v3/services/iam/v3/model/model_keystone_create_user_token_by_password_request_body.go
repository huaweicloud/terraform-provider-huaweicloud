package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneCreateUserTokenByPasswordRequestBody struct {
	Auth *PwdAuth `json:"auth"`
}

func (o KeystoneCreateUserTokenByPasswordRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateUserTokenByPasswordRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneCreateUserTokenByPasswordRequestBody", string(data)}, " ")
}
