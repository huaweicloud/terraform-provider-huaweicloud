package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateNewTaskResponse Response Object
type CreateNewTaskResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 任务id
	TaskId         *int32 `json:"task_id,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateNewTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateNewTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateNewTaskResponse", string(data)}, " ")
}
