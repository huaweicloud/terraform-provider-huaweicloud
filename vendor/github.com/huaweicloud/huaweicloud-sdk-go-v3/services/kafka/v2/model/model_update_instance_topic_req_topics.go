package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceTopicReqTopics 修改的topic。
type UpdateInstanceTopicReqTopics struct {

	// topic名称，不支持修改。
	Id string `json:"id"`

	// 老化时间，单位小时。
	RetentionTime *int32 `json:"retention_time,omitempty"`

	// 是否同步复制。
	SyncReplication *bool `json:"sync_replication,omitempty"`

	// 是否同步落盘。
	SyncMessageFlush *bool `json:"sync_message_flush,omitempty"`

	// 分区数。
	NewPartitionNumbers *int32 `json:"new_partition_numbers,omitempty"`

	// 增加分区时指定broker列表
	NewPartitionBrokers *[]int32 `json:"new_partition_brokers,omitempty"`

	// topic配置
	TopicOtherConfigs *[]CreateInstanceTopicReqTopicOtherConfigs `json:"topic_other_configs,omitempty"`

	// topic描述
	TopicDesc *string `json:"topic_desc,omitempty"`
}

func (o UpdateInstanceTopicReqTopics) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceTopicReqTopics struct{}"
	}

	return strings.Join([]string{"UpdateInstanceTopicReqTopics", string(data)}, " ")
}
