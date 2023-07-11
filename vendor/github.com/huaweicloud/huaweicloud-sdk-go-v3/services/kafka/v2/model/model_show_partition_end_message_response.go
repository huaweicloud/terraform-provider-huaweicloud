package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPartitionEndMessageResponse Response Object
type ShowPartitionEndMessageResponse struct {

	// Topic名称。
	Topic *string `json:"topic,omitempty"`

	// 分区编号。
	Partition *int32 `json:"partition,omitempty"`

	// 最新消息位置。
	Offset *int32 `json:"offset,omitempty"`

	// 生产消息的时间。 格式为Unix时间戳。单位为毫秒。
	Timestamp      *int64 `json:"timestamp,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowPartitionEndMessageResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPartitionEndMessageResponse struct{}"
	}

	return strings.Join([]string{"ShowPartitionEndMessageResponse", string(data)}, " ")
}
