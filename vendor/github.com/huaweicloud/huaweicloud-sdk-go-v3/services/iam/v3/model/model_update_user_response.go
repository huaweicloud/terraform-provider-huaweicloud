package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateUserResponse struct {
	User           *UpdateUserResult `json:"user,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o UpdateUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateUserResponse struct{}"
	}

	return strings.Join([]string{"UpdateUserResponse", string(data)}, " ")
}
