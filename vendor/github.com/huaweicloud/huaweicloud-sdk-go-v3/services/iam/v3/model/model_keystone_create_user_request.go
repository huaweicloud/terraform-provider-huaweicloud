package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneCreateUserRequest struct {
	Body *KeystoneCreateUserRequestBody `json:"body,omitempty"`
}

func (o KeystoneCreateUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateUserRequest struct{}"
	}

	return strings.Join([]string{"KeystoneCreateUserRequest", string(data)}, " ")
}
