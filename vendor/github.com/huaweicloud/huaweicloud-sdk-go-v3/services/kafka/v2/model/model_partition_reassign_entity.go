package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PartitionReassignEntity struct {

	// topic名称
	Topic string `json:"topic"`

	// 分区重平衡到的broker列表，自动生成分配方案时需指定该参数。
	Brokers *[]int32 `json:"brokers,omitempty"`

	// 副本因子，自动生成分配方案时可指定。
	ReplicationFactor *int32 `json:"replication_factor,omitempty"`

	// 手动指定的分配方案。brokers参数与该参数不能同时为空。
	Assignment *[]TopicAssignment `json:"assignment,omitempty"`
}

func (o PartitionReassignEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionReassignEntity struct{}"
	}

	return strings.Join([]string{"PartitionReassignEntity", string(data)}, " ")
}
