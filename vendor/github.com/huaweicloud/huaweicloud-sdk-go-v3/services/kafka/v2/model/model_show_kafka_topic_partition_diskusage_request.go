package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowKafkaTopicPartitionDiskusageRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 占用磁盘大小，默认值1G (1K ,1M , 1G)。
	MinSize *string `json:"minSize,omitempty"`

	// 占用磁盘大小，查询top N。
	Top *string `json:"top,omitempty"`

	// 占用磁盘大小，查询大于占比的分区。
	Percentage *string `json:"percentage,omitempty"`
}

func (o ShowKafkaTopicPartitionDiskusageRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowKafkaTopicPartitionDiskusageRequest struct{}"
	}

	return strings.Join([]string{"ShowKafkaTopicPartitionDiskusageRequest", string(data)}, " ")
}
