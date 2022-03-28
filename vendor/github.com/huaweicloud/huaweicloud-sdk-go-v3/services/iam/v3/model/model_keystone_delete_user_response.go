package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneDeleteUserResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneDeleteUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneDeleteUserResponse struct{}"
	}

	return strings.Join([]string{"KeystoneDeleteUserResponse", string(data)}, " ")
}
