package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ScaleNodePoolSpec 伸缩节点池请求详细参数
type ScaleNodePoolSpec struct {

	// 节点池期望节点数
	DesiredNodeCount int32 `json:"desiredNodeCount"`

	// 扩缩容的节点池，只能填一个伸缩组，如果要伸缩默认伸缩组填default
	ScaleGroups []string `json:"scaleGroups"`

	Options *ScaleNodePoolOptions `json:"options,omitempty"`
}

func (o ScaleNodePoolSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleNodePoolSpec struct{}"
	}

	return strings.Join([]string{"ScaleNodePoolSpec", string(data)}, " ")
}
