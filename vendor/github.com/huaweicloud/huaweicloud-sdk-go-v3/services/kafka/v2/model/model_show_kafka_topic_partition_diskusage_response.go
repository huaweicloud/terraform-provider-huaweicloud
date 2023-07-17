package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowKafkaTopicPartitionDiskusageResponse Response Object
type ShowKafkaTopicPartitionDiskusageResponse struct {

	// Broker列表。
	BrokerList     *[]DiskusageEntity `json:"broker_list,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ShowKafkaTopicPartitionDiskusageResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowKafkaTopicPartitionDiskusageResponse struct{}"
	}

	return strings.Join([]string{"ShowKafkaTopicPartitionDiskusageResponse", string(data)}, " ")
}
