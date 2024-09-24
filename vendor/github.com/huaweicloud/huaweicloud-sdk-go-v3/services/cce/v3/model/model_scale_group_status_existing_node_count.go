package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ScaleGroupStatusExistingNodeCount 伸缩组存量节点统计信息
type ScaleGroupStatusExistingNodeCount struct {

	// 按需计费节点个数
	PostPaid *int32 `json:"postPaid,omitempty"`

	// 包年包月节点个数
	PrePaid *int32 `json:"prePaid,omitempty"`

	// 按需计费和包年包月节点总数
	Total *int32 `json:"total,omitempty"`
}

func (o ScaleGroupStatusExistingNodeCount) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleGroupStatusExistingNodeCount struct{}"
	}

	return strings.Join([]string{"ScaleGroupStatusExistingNodeCount", string(data)}, " ")
}
