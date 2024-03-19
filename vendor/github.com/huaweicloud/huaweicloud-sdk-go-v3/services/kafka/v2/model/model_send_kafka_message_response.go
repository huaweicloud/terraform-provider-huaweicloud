package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SendKafkaMessageResponse Response Object
type SendKafkaMessageResponse struct {

	// Kafka的topic
	Topic *string `json:"topic,omitempty"`

	// 消息内容
	Body *string `json:"body,omitempty"`

	// topic的分区信息等
	PropertyList   *[]interface{} `json:"property_list,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o SendKafkaMessageResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SendKafkaMessageResponse struct{}"
	}

	return strings.Join([]string{"SendKafkaMessageResponse", string(data)}, " ")
}
