package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 当前规格实例的属性。
type ListEnginePropertiesEntity struct {

	// 每个Broker的最大分区数。
	MaxPartitionPerBroker *string `json:"max_partition_per_broker,omitempty"`

	// Broker的最大个数。
	MaxBroker *string `json:"max_broker,omitempty"`

	// 每个节点的最大存储。单位为GB。
	MaxStoragePerNode *string `json:"max_storage_per_node,omitempty"`

	// 每个Broker的最大消费者数。
	MaxConsumerPerBroker *string `json:"max_consumer_per_broker,omitempty"`

	// Broker的最小个数。
	MinBroker *string `json:"min_broker,omitempty"`

	// 每个Broker的最大带宽。
	MaxBandwidthPerBroker *string `json:"max_bandwidth_per_broker,omitempty"`

	// 每个节点的最小存储。单位为GB。
	MinStoragePerNode *string `json:"min_storage_per_node,omitempty"`

	// 每个Broker的最大TPS。
	MaxTpsPerBroker *string `json:"max_tps_per_broker,omitempty"`

	// product_id的别名。
	ProductAlias *string `json:"product_alias,omitempty"`
}

func (o ListEnginePropertiesEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEnginePropertiesEntity struct{}"
	}

	return strings.Join([]string{"ListEnginePropertiesEntity", string(data)}, " ")
}
