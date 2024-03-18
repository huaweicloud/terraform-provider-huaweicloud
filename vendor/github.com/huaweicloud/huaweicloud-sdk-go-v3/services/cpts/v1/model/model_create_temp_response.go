package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTempResponse Response Object
type CreateTempResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 事务id
	TempId *int32 `json:"tempId,omitempty"`

	// 响应消息
	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateTempResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTempResponse struct{}"
	}

	return strings.Join([]string{"CreateTempResponse", string(data)}, " ")
}
