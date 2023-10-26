package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type KafkaTopicProducerResponseProducers struct {

	// 生产者地址
	ProducerAddress *string `json:"producer_address,omitempty"`

	// broker地址
	BrokerAddress *string `json:"broker_address,omitempty"`

	// 加入时间
	JoinTime *int64 `json:"join_time,omitempty"`
}

func (o KafkaTopicProducerResponseProducers) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KafkaTopicProducerResponseProducers struct{}"
	}

	return strings.Join([]string{"KafkaTopicProducerResponseProducers", string(data)}, " ")
}
