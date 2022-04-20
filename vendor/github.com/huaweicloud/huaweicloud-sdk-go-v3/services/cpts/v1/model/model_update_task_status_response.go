package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateTaskStatusResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// message

	Message *string `json:"message,omitempty"`
	// extend

	Extend *string `json:"extend,omitempty"`

	Result         *UpdateTaskStatusResult `json:"result,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o UpdateTaskStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskStatusResponse struct{}"
	}

	return strings.Join([]string{"UpdateTaskStatusResponse", string(data)}, " ")
}
