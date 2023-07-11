package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateNewCaseResponse Response Object
type CreateNewCaseResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	Json *CreaseCaseResponseJson `json:"json,omitempty"`

	// 响应消息
	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateNewCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateNewCaseResponse struct{}"
	}

	return strings.Join([]string{"CreateNewCaseResponse", string(data)}, " ")
}
