package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowBackgroundTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 任务ID。
	TaskId string `json:"task_id"`
}

func (o ShowBackgroundTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBackgroundTaskRequest struct{}"
	}

	return strings.Join([]string{"ShowBackgroundTaskRequest", string(data)}, " ")
}
