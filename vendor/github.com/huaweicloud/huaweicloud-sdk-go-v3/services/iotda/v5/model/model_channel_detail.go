package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 物联网平台转发数据的通道配置参数。
type ChannelDetail struct {
	HttpForwarding *HttpForwarding `json:"http_forwarding,omitempty"`

	DisForwarding *DisForwarding `json:"dis_forwarding,omitempty"`

	ObsForwarding *ObsForwarding `json:"obs_forwarding,omitempty"`

	AmqpForwarding *AmqpForwarding `json:"amqp_forwarding,omitempty"`

	DmsKafkaForwarding *DmsKafkaForwarding `json:"dms_kafka_forwarding,omitempty"`
}

func (o ChannelDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChannelDetail struct{}"
	}

	return strings.Join([]string{"ChannelDetail", string(data)}, " ")
}
