package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateKafkaUserClientQuotaTaskRequest Request Object
type UpdateKafkaUserClientQuotaTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *UpdateKafkaUserClientQuotaTaskReq `json:"body,omitempty"`
}

func (o UpdateKafkaUserClientQuotaTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateKafkaUserClientQuotaTaskRequest struct{}"
	}

	return strings.Join([]string{"UpdateKafkaUserClientQuotaTaskRequest", string(data)}, " ")
}
