package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChannelDetail 物联网平台转发数据的通道配置参数。
type ChannelDetail struct {
	HttpForwarding *HttpForwarding `json:"http_forwarding,omitempty"`

	DisForwarding *DisForwarding `json:"dis_forwarding,omitempty"`

	ObsForwarding *ObsForwarding `json:"obs_forwarding,omitempty"`

	AmqpForwarding *AmqpForwarding `json:"amqp_forwarding,omitempty"`

	DmsKafkaForwarding *DmsKafkaForwarding `json:"dms_kafka_forwarding,omitempty"`

	RomaForwarding *RomaForwarding `json:"roma_forwarding,omitempty"`

	MysqlForwarding *MysqlForwarding `json:"mysql_forwarding,omitempty"`

	InfluxdbForwarding *InfluxDbForwarding `json:"influxdb_forwarding,omitempty"`

	FunctiongraphForwarding *FunctionGraphForwarding `json:"functiongraph_forwarding,omitempty"`

	MrsKafkaForwarding *MrsKafkaForwarding `json:"mrs_kafka_forwarding,omitempty"`

	DmsRocketmqForwarding *DmsRocketMqForwarding `json:"dms_rocketmq_forwarding,omitempty"`
}

func (o ChannelDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChannelDetail struct{}"
	}

	return strings.Join([]string{"ChannelDetail", string(data)}, " ")
}
