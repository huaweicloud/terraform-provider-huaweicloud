package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ExtendProductInfoEntity struct {

	// 实例类型。
	Type *string `json:"type,omitempty"`

	// 产品ID。
	ProductId *string `json:"product_id,omitempty"`

	// 该产品使用的ECS规格。
	EcsFlavorId *string `json:"ecs_flavor_id,omitempty"`

	// 支持的CPU架构类型。
	ArchTypes *[]string `json:"arch_types,omitempty"`

	// 支持的计费模式类型。
	ChargingMode *[]string `json:"charging_mode,omitempty"`

	// 磁盘IO信息。
	Ios *[]ExtendProductIosEntity `json:"ios,omitempty"`

	// 支持的特性功能。
	SupportFeatures *[]ExtendProductSupportFeaturesEntity `json:"support_features,omitempty"`

	Properties *ExtendProductPropertiesEntity `json:"properties,omitempty"`

	// 有可用资源的可用区列表。
	AvailableZones *[]string `json:"available_zones,omitempty"`

	// 资源售罄的可用区列表。
	UnavailableZones *[]string `json:"unavailable_zones,omitempty"`
}

func (o ExtendProductInfoEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExtendProductInfoEntity struct{}"
	}

	return strings.Join([]string{"ExtendProductInfoEntity", string(data)}, " ")
}
