package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateTaskRequest struct {
	// 任务id

	TaskId int32 `json:"task_id"`

	Body *UpdateTaskRequestBody `json:"body,omitempty"`
}

func (o UpdateTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskRequest struct{}"
	}

	return strings.Join([]string{"UpdateTaskRequest", string(data)}, " ")
}
