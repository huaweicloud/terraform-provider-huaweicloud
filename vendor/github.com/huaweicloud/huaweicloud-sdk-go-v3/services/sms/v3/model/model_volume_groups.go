package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 逻辑卷组信息
type VolumeGroups struct {
	// Pv信息

	Components *string `json:"components,omitempty"`
	// 剩余空间

	FreeSize *int64 `json:"free_size,omitempty"`
	// lv信息

	LogicalVolumes *[]LogicalVolumes `json:"logical_volumes,omitempty"`
	// 名称

	Name *string `json:"name,omitempty"`
	// 大小

	Size *int64 `json:"size,omitempty"`
}

func (o VolumeGroups) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VolumeGroups struct{}"
	}

	return strings.Join([]string{"VolumeGroups", string(data)}, " ")
}
