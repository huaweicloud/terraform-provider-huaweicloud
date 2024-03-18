package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowMessagesRespMessages struct {

	// topic名称。
	Topic *string `json:"topic,omitempty"`

	// 分区编号。
	Partition *int32 `json:"partition,omitempty"`

	// 消息编号。
	MessageOffset *int64 `json:"message_offset,omitempty"`

	// 消息大小，单位字节。
	Size *int32 `json:"size,omitempty"`

	// 生产消息的时间。 格式为Unix时间戳。单位为毫秒。
	Timestamp *int64 `json:"timestamp,omitempty"`
}

func (o ShowMessagesRespMessages) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMessagesRespMessages struct{}"
	}

	return strings.Join([]string{"ShowMessagesRespMessages", string(data)}, " ")
}
