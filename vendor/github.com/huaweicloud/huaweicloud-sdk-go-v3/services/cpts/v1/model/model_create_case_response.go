package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateCaseResponse struct {
	// code

	Code *string `json:"code,omitempty"`

	Json *CreateCaseResultJson `json:"json,omitempty"`
	// message

	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCaseResponse struct{}"
	}

	return strings.Join([]string{"CreateCaseResponse", string(data)}, " ")
}
