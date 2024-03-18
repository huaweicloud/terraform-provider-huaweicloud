package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SendKafkaMessageRequestBody Kafka生产消息请求体
type SendKafkaMessageRequestBody struct {

	// Kafka的topic
	Topic string `json:"topic"`

	// 消息内容
	Body string `json:"body"`

	// topic的分区信息等
	PropertyList []SendKafkaMessageRequestBodyPropertyList `json:"property_list"`
}

func (o SendKafkaMessageRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SendKafkaMessageRequestBody struct{}"
	}

	return strings.Join([]string{"SendKafkaMessageRequestBody", string(data)}, " ")
}
