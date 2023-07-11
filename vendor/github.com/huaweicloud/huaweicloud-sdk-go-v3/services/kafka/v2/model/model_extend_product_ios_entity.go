package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExtendProductIosEntity 磁盘IO信息。
type ExtendProductIosEntity struct {

	// 存储IO规格。
	IoSpec *string `json:"io_spec,omitempty"`

	// 有可用资源的可用区列表。
	AvailableZones *[]string `json:"available_zones,omitempty"`

	// IO类型。
	Type *string `json:"type,omitempty"`

	// 资源售罄的可用区列表。
	UnavailableZones *[]string `json:"unavailable_zones,omitempty"`
}

func (o ExtendProductIosEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExtendProductIosEntity struct{}"
	}

	return strings.Join([]string{"ExtendProductIosEntity", string(data)}, " ")
}
