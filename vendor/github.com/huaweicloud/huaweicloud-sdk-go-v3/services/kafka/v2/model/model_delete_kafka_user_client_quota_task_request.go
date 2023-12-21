package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteKafkaUserClientQuotaTaskRequest Request Object
type DeleteKafkaUserClientQuotaTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *DeleteKafkaUserClientQuotaTaskReq `json:"body,omitempty"`
}

func (o DeleteKafkaUserClientQuotaTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteKafkaUserClientQuotaTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteKafkaUserClientQuotaTaskRequest", string(data)}, " ")
}
