package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateKafkaConsumerGroupRequest Request Object
type CreateKafkaConsumerGroupRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateGroupReq `json:"body,omitempty"`
}

func (o CreateKafkaConsumerGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateKafkaConsumerGroupRequest struct{}"
	}

	return strings.Join([]string{"CreateKafkaConsumerGroupRequest", string(data)}, " ")
}
