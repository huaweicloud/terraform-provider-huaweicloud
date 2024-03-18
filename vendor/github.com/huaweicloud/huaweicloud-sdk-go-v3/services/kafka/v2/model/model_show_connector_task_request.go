package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowConnectorTaskRequest Request Object
type ShowConnectorTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Smart Connector任务ID。
	TaskId string `json:"task_id"`
}

func (o ShowConnectorTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowConnectorTaskRequest struct{}"
	}

	return strings.Join([]string{"ShowConnectorTaskRequest", string(data)}, " ")
}
