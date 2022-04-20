package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteTaskRequest struct {
	// 任务id

	TaskId int32 `json:"task_id"`
}

func (o DeleteTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteTaskRequest", string(data)}, " ")
}
