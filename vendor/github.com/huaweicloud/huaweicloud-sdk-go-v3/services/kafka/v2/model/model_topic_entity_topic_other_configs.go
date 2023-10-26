package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TopicEntityTopicOtherConfigs struct {

	// 配置名称
	Name *string `json:"name,omitempty"`

	// 配置有效值
	ValidValues *string `json:"valid_values,omitempty"`

	// 配置默认值
	DefaultValue *string `json:"default_value,omitempty"`

	// 配置类型：dynamic/static
	ConfigType *string `json:"config_type,omitempty"`

	// 配置值
	Value *string `json:"value,omitempty"`

	// 配置值类型
	ValueType *string `json:"value_type,omitempty"`
}

func (o TopicEntityTopicOtherConfigs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TopicEntityTopicOtherConfigs struct{}"
	}

	return strings.Join([]string{"TopicEntityTopicOtherConfigs", string(data)}, " ")
}
