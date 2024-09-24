package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodeInfo struct {

	// 节点类型。
	NodeType *string `json:"node_type,omitempty"`

	// 节点id。
	NodeId *string `json:"node_id,omitempty"`
}

func (o NodeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeInfo struct{}"
	}

	return strings.Join([]string{"NodeInfo", string(data)}, " ")
}
