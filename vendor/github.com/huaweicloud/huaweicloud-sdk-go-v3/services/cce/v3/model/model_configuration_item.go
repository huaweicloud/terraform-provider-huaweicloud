package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ConfigurationItem struct {

	// 组件配置项名称
	Name *string `json:"name,omitempty"`

	// 组件配置项值
	Value *interface{} `json:"value,omitempty"`
}

func (o ConfigurationItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigurationItem struct{}"
	}

	return strings.Join([]string{"ConfigurationItem", string(data)}, " ")
}
