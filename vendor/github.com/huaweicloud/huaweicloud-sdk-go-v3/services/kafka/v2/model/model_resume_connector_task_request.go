package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResumeConnectorTaskRequest Request Object
type ResumeConnectorTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Smart Connect任务ID。
	TaskId string `json:"task_id"`
}

func (o ResumeConnectorTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResumeConnectorTaskRequest struct{}"
	}

	return strings.Join([]string{"ResumeConnectorTaskRequest", string(data)}, " ")
}
