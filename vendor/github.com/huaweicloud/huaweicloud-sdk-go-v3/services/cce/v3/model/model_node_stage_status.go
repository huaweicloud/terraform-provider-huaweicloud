package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeStageStatus 节点检查状态
type NodeStageStatus struct {
	NodeInfo *NodeInfo `json:"nodeInfo,omitempty"`

	// 检查项状态集合
	ItemsStatus *[]PreCheckItemStatus `json:"itemsStatus,omitempty"`
}

func (o NodeStageStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeStageStatus struct{}"
	}

	return strings.Join([]string{"NodeStageStatus", string(data)}, " ")
}
