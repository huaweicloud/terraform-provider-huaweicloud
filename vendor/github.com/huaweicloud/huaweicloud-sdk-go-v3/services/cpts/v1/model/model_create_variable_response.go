package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateVariableResponse struct {
	// code

	Code *string `json:"code,omitempty"`

	Json *CreateVariableResultJson `json:"json,omitempty"`
	// message

	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateVariableResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateVariableResponse struct{}"
	}

	return strings.Join([]string{"CreateVariableResponse", string(data)}, " ")
}
