package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowInstanceExtendProductInfoRespValues struct {

	// 规格详情。
	Detail *[]ShowInstanceExtendProductInfoRespDetail `json:"detail,omitempty"`

	// 实例类型。
	Name *string `json:"name,omitempty"`

	// 资源售罄的可用区列表。
	UnavailableZones *[]string `json:"unavailable_zones,omitempty"`

	// 有可用资源的可用区列表。
	AvailableZones *[]string `json:"available_zones,omitempty"`
}

func (o ShowInstanceExtendProductInfoRespValues) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceExtendProductInfoRespValues struct{}"
	}

	return strings.Join([]string{"ShowInstanceExtendProductInfoRespValues", string(data)}, " ")
}
