package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateKafkaConsumerGroupResponse Response Object
type CreateKafkaConsumerGroupResponse struct {

	// 创建结果
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateKafkaConsumerGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateKafkaConsumerGroupResponse struct{}"
	}

	return strings.Join([]string{"CreateKafkaConsumerGroupResponse", string(data)}, " ")
}
