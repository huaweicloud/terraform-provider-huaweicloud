package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RestartConnectorTaskRequest Request Object
type RestartConnectorTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Smart Connect任务ID。
	TaskId string `json:"task_id"`
}

func (o RestartConnectorTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestartConnectorTaskRequest struct{}"
	}

	return strings.Join([]string{"RestartConnectorTaskRequest", string(data)}, " ")
}
