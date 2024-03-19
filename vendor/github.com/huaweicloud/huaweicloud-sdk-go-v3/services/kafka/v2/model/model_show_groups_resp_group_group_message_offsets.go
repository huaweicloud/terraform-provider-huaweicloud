package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowGroupsRespGroupGroupMessageOffsets struct {

	// 分区编号。
	Partition *int32 `json:"partition,omitempty"`

	// 剩余可消费消息数，即消息堆积数。
	Lag *int64 `json:"lag,omitempty"`

	// topic名称。
	Topic *string `json:"topic,omitempty"`

	// 当前消费进度。
	MessageCurrentOffset *int64 `json:"message_current_offset,omitempty"`

	// 最大消息位置（LEO）。
	MessageLogEndOffset *int64 `json:"message_log_end_offset,omitempty"`
}

func (o ShowGroupsRespGroupGroupMessageOffsets) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGroupsRespGroupGroupMessageOffsets struct{}"
	}

	return strings.Join([]string{"ShowGroupsRespGroupGroupMessageOffsets", string(data)}, " ")
}
