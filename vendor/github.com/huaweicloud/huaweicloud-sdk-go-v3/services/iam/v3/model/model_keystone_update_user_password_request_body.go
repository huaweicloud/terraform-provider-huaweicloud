package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneUpdateUserPasswordRequestBody struct {
	User *KeystoneUpdatePasswordOption `json:"user"`
}

func (o KeystoneUpdateUserPasswordRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateUserPasswordRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateUserPasswordRequestBody", string(data)}, " ")
}
