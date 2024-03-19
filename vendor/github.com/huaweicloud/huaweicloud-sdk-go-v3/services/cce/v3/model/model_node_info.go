package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeInfo 节点检查状态
type NodeInfo struct {

	// 节点UID
	Uid *string `json:"uid,omitempty"`

	// 节点名称
	Name *string `json:"name,omitempty"`

	// 状态
	Status *string `json:"status,omitempty"`

	// 节点类型
	NodeType *string `json:"nodeType,omitempty"`
}

func (o NodeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeInfo struct{}"
	}

	return strings.Join([]string{"NodeInfo", string(data)}, " ")
}
