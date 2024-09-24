package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AutopilotConfigurationItem struct {

	// 组件配置项名称
	Name *string `json:"name,omitempty"`

	// 组件配置项值
	Value *interface{} `json:"value,omitempty"`
}

func (o AutopilotConfigurationItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotConfigurationItem struct{}"
	}

	return strings.Join([]string{"AutopilotConfigurationItem", string(data)}, " ")
}
