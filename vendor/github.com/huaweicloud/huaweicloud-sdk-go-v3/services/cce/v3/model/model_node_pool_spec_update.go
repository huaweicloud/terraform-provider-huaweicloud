package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodePoolSpecUpdate
type NodePoolSpecUpdate struct {
	NodeTemplate *NodeSpecUpdate `json:"nodeTemplate"`

	// 节点池初始化节点个数。查询时为节点池目标节点数量。默认值为0。
	InitialNodeCount int32 `json:"initialNodeCount"`

	Autoscaling *NodePoolNodeAutoscaling `json:"autoscaling"`
}

func (o NodePoolSpecUpdate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodePoolSpecUpdate struct{}"
	}

	return strings.Join([]string{"NodePoolSpecUpdate", string(data)}, " ")
}
