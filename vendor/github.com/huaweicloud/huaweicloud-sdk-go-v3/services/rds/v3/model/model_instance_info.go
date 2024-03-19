package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InstanceInfo 数据库实例信息。
type InstanceInfo struct {

	// 数据库主实例或只读实例ID。
	Id *string `json:"id,omitempty"`

	// 节点状态。
	Status *string `json:"status,omitempty"`

	// 数据库实例名称。
	Name *string `json:"name,omitempty"`

	// 数据库实例读权重。
	Weight *int32 `json:"weight,omitempty"`

	// 可用区信息。
	AvailableZones *[]MysqlAvailableZoneInfo `json:"available_zones,omitempty"`
}

func (o InstanceInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceInfo struct{}"
	}

	return strings.Join([]string{"InstanceInfo", string(data)}, " ")
}
