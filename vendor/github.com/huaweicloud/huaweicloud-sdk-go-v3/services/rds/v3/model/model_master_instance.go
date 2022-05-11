package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MasterInstance struct {

	// 主实例ID。
	Id string `json:"id"`

	// 主实例状态。
	Status string `json:"status"`

	// 主实例名称。
	Name string `json:"name"`

	// 主实例读写分离权重。
	Weight int32 `json:"weight"`

	// 可用区信息。
	AvailableZones []AvailableZone `json:"available_zones"`

	// 主实例CPU个数。
	CpuNum int32 `json:"cpu_num"`
}

func (o MasterInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MasterInstance struct{}"
	}

	return strings.Join([]string{"MasterInstance", string(data)}, " ")
}
