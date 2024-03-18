package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResetMessageOffsetReq struct {

	// topic名称。
	Topic *string `json:"topic,omitempty"`

	// 分区编号，默认值为-1，若传入值为-1，则重置所有分区。
	Partition int32 `json:"partition"`

	// 重置消费进度到指定偏移量。 - 如果传入offset小于当前最小的offset，则重置到最小的offset。 - 如果大于最大的offset，则重置到最大的offset。  **message_offset、timestamp二者必选其一。**
	MessageOffset *int64 `json:"message_offset,omitempty"`

	// 重置消费进度到指定时间，格式为unix时间戳，单位为毫秒。 - 如果传入timestamp早于当前最早的timestamp，则重置到最早的timestamp。 - 如果晚于最晚的timestamp，则重置到最晚的timestamp。  **message_offset、timestamp二者必选其一。**
	Timestamp *int64 `json:"timestamp,omitempty"`
}

func (o ResetMessageOffsetReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetMessageOffsetReq struct{}"
	}

	return strings.Join([]string{"ResetMessageOffsetReq", string(data)}, " ")
}
