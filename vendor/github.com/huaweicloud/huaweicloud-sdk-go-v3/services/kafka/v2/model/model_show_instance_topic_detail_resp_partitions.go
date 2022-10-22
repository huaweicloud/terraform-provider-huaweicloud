package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowInstanceTopicDetailRespPartitions struct {

	// 分区ID。
	Partition *int32 `json:"partition,omitempty"`

	// leader副本所在节点的id。
	Leader *int32 `json:"leader,omitempty"`

	// 分区leader副本的LEO（Log End Offset）。
	Leo *int32 `json:"leo,omitempty"`

	// 分区高水位（HW，High Watermark）。
	Hw *int32 `json:"hw,omitempty"`

	// 分区leader副本的LSO（Log Start Offset）。
	Lso *int32 `json:"lso,omitempty"`

	// 分区上次写入消息的时间。  格式为Unix时间戳。  单位：毫秒。
	LastUpdateTimestamp *int64 `json:"last_update_timestamp,omitempty"`

	// 副本列表。
	Replicas *[]ShowInstanceTopicDetailRespReplicas `json:"replicas,omitempty"`
}

func (o ShowInstanceTopicDetailRespPartitions) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceTopicDetailRespPartitions struct{}"
	}

	return strings.Join([]string{"ShowInstanceTopicDetailRespPartitions", string(data)}, " ")
}
