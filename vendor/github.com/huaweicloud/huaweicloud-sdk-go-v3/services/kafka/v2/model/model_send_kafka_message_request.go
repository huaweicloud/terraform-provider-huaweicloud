package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SendKafkaMessageRequest Request Object
type SendKafkaMessageRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 动作ID，生产消息对应的action_id为send。
	ActionId string `json:"action_id"`

	Body *SendKafkaMessageRequestBody `json:"body,omitempty"`
}

func (o SendKafkaMessageRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SendKafkaMessageRequest struct{}"
	}

	return strings.Join([]string{"SendKafkaMessageRequest", string(data)}, " ")
}
