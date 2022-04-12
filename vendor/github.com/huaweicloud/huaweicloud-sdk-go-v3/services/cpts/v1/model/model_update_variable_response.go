package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateVariableResponse struct {
	// code

	Code *string `json:"code,omitempty"`

	Json *CreateVariableResultJson `json:"json,omitempty"`
	// message

	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateVariableResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateVariableResponse struct{}"
	}

	return strings.Join([]string{"UpdateVariableResponse", string(data)}, " ")
}
