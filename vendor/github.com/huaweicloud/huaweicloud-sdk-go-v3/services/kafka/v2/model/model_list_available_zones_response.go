package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAvailableZonesResponse Response Object
type ListAvailableZonesResponse struct {

	// 区域ID。
	RegionId *string `json:"region_id,omitempty"`

	// 可用区数组。
	AvailableZones *[]AvailableZonesResp `json:"available_zones,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o ListAvailableZonesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAvailableZonesResponse struct{}"
	}

	return strings.Join([]string{"ListAvailableZonesResponse", string(data)}, " ")
}
