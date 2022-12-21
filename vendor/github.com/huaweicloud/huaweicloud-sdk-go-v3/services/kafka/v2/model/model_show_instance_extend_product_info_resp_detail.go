package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowInstanceExtendProductInfoRespDetail struct {

	// 单位时间内的消息量最大值。
	Tps *string `json:"tps,omitempty"`

	// 消息存储空间。
	Storage *string `json:"storage,omitempty"`

	// Kafka实例的分区数量。
	PartitionNum *string `json:"partition_num,omitempty"`

	// 产品ID。
	ProductId *string `json:"product_id,omitempty"`

	// 规格ID。
	SpecCode *string `json:"spec_code,omitempty"`

	// IO信息。
	Io *[]ListProductsRespIo `json:"io,omitempty"`

	// Kafka实例的基准带宽。
	Bandwidth *string `json:"bandwidth,omitempty"`

	// Kafka实例最大消费组数参考值。
	RecommendMaxConsGroups *string `json:"recommend_max_consGroups,omitempty"`

	// 资源售罄的可用区列表。
	UnavailableZones *[]string `json:"unavailable_zones,omitempty"`

	// 有可用资源的可用区列表。
	AvailableZones *[]string `json:"available_zones,omitempty"`

	// 该产品规格对应的虚拟机规格。
	EcsFlavorId *string `json:"ecs_flavor_id,omitempty"`

	// 实例规格架构类型。当前仅支持X86。
	ArchType *string `json:"arch_type,omitempty"`
}

func (o ShowInstanceExtendProductInfoRespDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceExtendProductInfoRespDetail struct{}"
	}

	return strings.Join([]string{"ShowInstanceExtendProductInfoRespDetail", string(data)}, " ")
}
