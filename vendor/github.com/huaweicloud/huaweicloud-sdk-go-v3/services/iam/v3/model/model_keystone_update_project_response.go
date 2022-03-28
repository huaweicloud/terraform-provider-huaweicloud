package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneUpdateProjectResponse struct {
	Project        *KeystoneUpdateProjectResult `json:"project,omitempty"`
	HttpStatusCode int                          `json:"-"`
}

func (o KeystoneUpdateProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateProjectResponse struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateProjectResponse", string(data)}, " ")
}
