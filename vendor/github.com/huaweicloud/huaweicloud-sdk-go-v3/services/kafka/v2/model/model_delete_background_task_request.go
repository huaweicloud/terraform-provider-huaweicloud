package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteBackgroundTaskRequest Request Object
type DeleteBackgroundTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 任务ID。
	TaskId string `json:"task_id"`
}

func (o DeleteBackgroundTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteBackgroundTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteBackgroundTaskRequest", string(data)}, " ")
}
