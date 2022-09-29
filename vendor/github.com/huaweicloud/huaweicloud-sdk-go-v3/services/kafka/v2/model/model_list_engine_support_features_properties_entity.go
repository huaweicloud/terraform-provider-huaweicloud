package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 实例支持的功能属性描述。
type ListEngineSupportFeaturesPropertiesEntity struct {

	// 转储功能的最大任务数。
	MaxTask *string `json:"max_task,omitempty"`

	// 转储功能的最小任务数。
	MinTask *string `json:"min_task,omitempty"`

	// 转储功能的最大节点数。
	MaxNode *string `json:"max_node,omitempty"`

	// 转储功能的最小节点数。
	MinNode *string `json:"min_node,omitempty"`
}

func (o ListEngineSupportFeaturesPropertiesEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEngineSupportFeaturesPropertiesEntity struct{}"
	}

	return strings.Join([]string{"ListEngineSupportFeaturesPropertiesEntity", string(data)}, " ")
}
