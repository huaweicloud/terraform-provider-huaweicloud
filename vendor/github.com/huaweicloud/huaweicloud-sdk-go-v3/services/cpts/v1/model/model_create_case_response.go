package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateCaseResponse Response Object
type CreateCaseResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	Json *CreateCaseResultJson `json:"json,omitempty"`

	// 响应消息
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
