package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteConnectorTaskRequest Request Object
type DeleteConnectorTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Smart Connector任务ID。
	TaskId string `json:"task_id"`
}

func (o DeleteConnectorTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConnectorTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteConnectorTaskRequest", string(data)}, " ")
}
