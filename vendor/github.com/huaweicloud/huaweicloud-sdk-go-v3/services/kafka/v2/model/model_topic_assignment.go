package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TopicAssignment struct {

	// 手动指定分配方案时的分区号。
	Partition *int32 `json:"partition,omitempty"`

	// 手动指定某个分区将要分配的broker列表
	PartitionBrokers *[]int32 `json:"partition_brokers,omitempty"`
}

func (o TopicAssignment) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TopicAssignment struct{}"
	}

	return strings.Join([]string{"TopicAssignment", string(data)}, " ")
}
