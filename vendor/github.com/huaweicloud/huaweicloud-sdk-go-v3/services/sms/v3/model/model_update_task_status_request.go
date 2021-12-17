package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateTaskStatusRequest struct {
	// 迁移任务ID

	TaskId string `json:"task_id"`

	Body *UpdateTaskStatusReq `json:"body,omitempty"`
}

func (o UpdateTaskStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskStatusRequest struct{}"
	}

	return strings.Join([]string{"UpdateTaskStatusRequest", string(data)}, " ")
}
