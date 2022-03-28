package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneCreateUserRequestBody struct {
	User *KeystoneCreateUserOption `json:"user"`
}

func (o KeystoneCreateUserRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateUserRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneCreateUserRequestBody", string(data)}, " ")
}
