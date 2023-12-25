package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateKafkaUserClientQuotaTaskRequest Request Object
type CreateKafkaUserClientQuotaTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateKafkaUserClientQuotaTaskReq `json:"body,omitempty"`
}

func (o CreateKafkaUserClientQuotaTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateKafkaUserClientQuotaTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateKafkaUserClientQuotaTaskRequest", string(data)}, " ")
}
