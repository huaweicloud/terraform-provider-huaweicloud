package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateDirectoryResponse Response Object
type CreateDirectoryResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 目录id
	DirectoryId *int32 `json:"directory_id,omitempty"`

	// 响应消息
	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateDirectoryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDirectoryResponse struct{}"
	}

	return strings.Join([]string{"CreateDirectoryResponse", string(data)}, " ")
}
