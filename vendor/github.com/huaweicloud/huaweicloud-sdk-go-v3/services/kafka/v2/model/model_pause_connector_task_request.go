package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PauseConnectorTaskRequest Request Object
type PauseConnectorTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Smart Connect任务ID。
	TaskId string `json:"task_id"`
}

func (o PauseConnectorTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PauseConnectorTaskRequest struct{}"
	}

	return strings.Join([]string{"PauseConnectorTaskRequest", string(data)}, " ")
}
