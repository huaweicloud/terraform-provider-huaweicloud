package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 支持的磁盘IO类型信息。
type ListEngineIosEntity struct {

	// 磁盘IO编码。
	IoSpec *string `json:"io_spec,omitempty"`

	// 磁盘类型。
	Type *string `json:"type,omitempty"`

	// 可用区。
	AvailableZones *[]string `json:"available_zones,omitempty"`

	// 不可用区。
	UnavailableZones *[]string `json:"unavailable_zones,omitempty"`
}

func (o ListEngineIosEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEngineIosEntity struct{}"
	}

	return strings.Join([]string{"ListEngineIosEntity", string(data)}, " ")
}
