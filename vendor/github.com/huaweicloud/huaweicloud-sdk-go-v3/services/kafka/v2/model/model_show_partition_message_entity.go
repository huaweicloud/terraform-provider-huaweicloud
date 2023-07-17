package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPartitionMessageEntity 消息体。
type ShowPartitionMessageEntity struct {

	// 消息的key。
	Key *string `json:"key,omitempty"`

	// 消息内容。
	Value *string `json:"value,omitempty"`

	// Topic名称。
	Topic *string `json:"topic,omitempty"`

	// 分区编号。
	Partition *int32 `json:"partition,omitempty"`

	// 消息位置。
	MessageOffset *int64 `json:"message_offset,omitempty"`

	// 消息大小，单位字节。
	Size *int32 `json:"size,omitempty"`

	// 生产消息的时间。 格式为Unix时间戳。单位为毫秒。
	Timestamp *int64 `json:"timestamp,omitempty"`
}

func (o ShowPartitionMessageEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPartitionMessageEntity struct{}"
	}

	return strings.Join([]string{"ShowPartitionMessageEntity", string(data)}, " ")
}
