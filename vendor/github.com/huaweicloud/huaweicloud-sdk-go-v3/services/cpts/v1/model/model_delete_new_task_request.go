package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteNewTaskRequest Request Object
type DeleteNewTaskRequest struct {

	// 任务id
	TaskId int32 `json:"task_id"`
}

func (o DeleteNewTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteNewTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteNewTaskRequest", string(data)}, " ")
}
