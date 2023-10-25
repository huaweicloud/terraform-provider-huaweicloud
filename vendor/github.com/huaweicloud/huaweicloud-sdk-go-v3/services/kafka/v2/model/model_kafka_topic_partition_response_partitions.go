package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type KafkaTopicPartitionResponsePartitions struct {

	// 分区ID
	Partition *int32 `json:"partition,omitempty"`

	// 起始偏移量
	StartOffset *int64 `json:"start_offset,omitempty"`

	// 最后偏移量
	LastOffset *int64 `json:"last_offset,omitempty"`

	// 分区消息数
	MessageCount *int64 `json:"message_count,omitempty"`

	// 最近更新时间
	LastUpdateTime *int64 `json:"last_update_time,omitempty"`
}

func (o KafkaTopicPartitionResponsePartitions) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KafkaTopicPartitionResponsePartitions struct{}"
	}

	return strings.Join([]string{"KafkaTopicPartitionResponsePartitions", string(data)}, " ")
}
