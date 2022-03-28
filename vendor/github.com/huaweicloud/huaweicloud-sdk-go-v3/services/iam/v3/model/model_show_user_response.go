package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowUserResponse struct {
	User           *ShowUserResult `json:"user,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ShowUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUserResponse struct{}"
	}

	return strings.Join([]string{"ShowUserResponse", string(data)}, " ")
}
