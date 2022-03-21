package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCreateUserResponse struct {
	User           *KeystoneCreateUserResult `json:"user,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o KeystoneCreateUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateUserResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCreateUserResponse", string(data)}, " ")
}
