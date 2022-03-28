package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateUserRequest struct {
	Body *CreateUserRequestBody `json:"body,omitempty"`
}

func (o CreateUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateUserRequest struct{}"
	}

	return strings.Join([]string{"CreateUserRequest", string(data)}, " ")
}
