package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateTaskSpeedRequest struct {
	// 主机迁移任务的id

	TaskId string `json:"task_id"`

	Body *UpdateTaskSpeedReq `json:"body,omitempty"`
}

func (o UpdateTaskSpeedRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskSpeedRequest struct{}"
	}

	return strings.Join([]string{"UpdateTaskSpeedRequest", string(data)}, " ")
}
