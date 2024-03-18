package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeRisks 节点风险来源
type NodeRisks struct {

	// 用户节点ID
	NodeID *string `json:"NodeID,omitempty"`
}

func (o NodeRisks) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeRisks struct{}"
	}

	return strings.Join([]string{"NodeRisks", string(data)}, " ")
}
