package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowKafkaUserClientQuotaRequest Request Object
type ShowKafkaUserClientQuotaRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 偏移量，表示查询该偏移量后面的记录。
	Offset *int32 `json:"offset,omitempty"`

	// 查询返回记录的数量限制。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ShowKafkaUserClientQuotaRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowKafkaUserClientQuotaRequest struct{}"
	}

	return strings.Join([]string{"ShowKafkaUserClientQuotaRequest", string(data)}, " ")
}
