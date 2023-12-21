package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodePriority 节点优先级批量配置
type NodePriority struct {
	NodeSelector *NodeSelector `json:"nodeSelector"`

	// 该批次节点的优先级，默认值为0，优先级最低，数值越大优先级越高
	Priority int32 `json:"priority"`
}

func (o NodePriority) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodePriority struct{}"
	}

	return strings.Join([]string{"NodePriority", string(data)}, " ")
}
