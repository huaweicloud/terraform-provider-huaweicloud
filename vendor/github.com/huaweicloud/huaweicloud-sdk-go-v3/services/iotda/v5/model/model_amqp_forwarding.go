package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// amqp queue配置信息
type AmqpForwarding struct {

	// **参数说明**：用于接收满足规则条件数据的amqp queue。
	QueueName string `json:"queue_name"`
}

func (o AmqpForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AmqpForwarding struct{}"
	}

	return strings.Join([]string{"AmqpForwarding", string(data)}, " ")
}
