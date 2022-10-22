package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DiskusageEntity struct {

	// Broker名称。
	BrokerName *string `json:"broker_name,omitempty"`

	// 磁盘容量。
	DataDiskSize *string `json:"data_disk_size,omitempty"`

	// 已使用的磁盘容量。
	DataDiskUse *string `json:"data_disk_use,omitempty"`

	// 剩余可用的磁盘容量。
	DataDiskFree *string `json:"data_disk_free,omitempty"`

	// 消息标签。
	DataDiskUsePercentage *string `json:"data_disk_use_percentage,omitempty"`

	// 消息标签。
	Status *string `json:"status,omitempty"`

	// topic磁盘容量使用列表。
	TopicList *[]DiskusageTopicEntity `json:"topic_list,omitempty"`
}

func (o DiskusageEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DiskusageEntity struct{}"
	}

	return strings.Join([]string{"DiskusageEntity", string(data)}, " ")
}
