package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTaskResponse Response Object
type CreateTaskResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 任务id
	TaskId         *int32 `json:"task_id,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateTaskResponse", string(data)}, " ")
}
