package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCreateProjectResponse struct {
	Project        *AuthProjectResult `json:"project,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o KeystoneCreateProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateProjectResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCreateProjectResponse", string(data)}, " ")
}
