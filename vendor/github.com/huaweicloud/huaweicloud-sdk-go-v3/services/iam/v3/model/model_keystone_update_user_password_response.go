package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneUpdateUserPasswordResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneUpdateUserPasswordResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateUserPasswordResponse struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateUserPasswordResponse", string(data)}, " ")
}
