package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlProxyNode struct {
	// 节点id。

	Id *string `json:"id,omitempty"`
	// 实例id。

	InstanceId *string `json:"instance_id,omitempty"`
	// 节点状态。

	Status *string `json:"status,omitempty"`
	// 节点名称。

	Name *string `json:"name,omitempty"`
	// 节点读写分离权重。

	Weight *int32 `json:"weight,omitempty"`
	// 可用区信息。

	AvailableZones *[]MysqlProxyAvailable `json:"available_zones,omitempty"`
}

func (o MysqlProxyNode) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlProxyNode struct{}"
	}

	return strings.Join([]string{"MysqlProxyNode", string(data)}, " ")
}
