package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AvailableZonesResp struct {

	// 是否售罄。
	SoldOut *bool `json:"soldOut,omitempty"`

	// 可用区ID。
	Id *string `json:"id,omitempty"`

	// 可用区编码。
	Code *string `json:"code,omitempty"`

	// 可用区名称。
	Name *string `json:"name,omitempty"`

	// 可用区端口号。
	Port *string `json:"port,omitempty"`

	// 可用区上是否还有可用资源。
	ResourceAvailability *string `json:"resource_availability,omitempty"`

	// 是否为默认可用区。
	DefaultAz *bool `json:"default_az,omitempty"`

	// 剩余时间。
	RemainTime *int64 `json:"remain_time,omitempty"`

	// 是否支持IPv6。
	Ipv6Enable *bool `json:"ipv6_enable,omitempty"`
}

func (o AvailableZonesResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AvailableZonesResp struct{}"
	}

	return strings.Join([]string{"AvailableZonesResp", string(data)}, " ")
}
