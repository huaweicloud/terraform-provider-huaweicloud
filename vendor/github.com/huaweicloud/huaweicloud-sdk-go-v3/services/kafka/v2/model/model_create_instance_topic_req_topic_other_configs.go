package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateInstanceTopicReqTopicOtherConfigs struct {

	// 配置名称
	Name *string `json:"name,omitempty"`

	// 配置值
	Value *string `json:"value,omitempty"`
}

func (o CreateInstanceTopicReqTopicOtherConfigs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceTopicReqTopicOtherConfigs struct{}"
	}

	return strings.Join([]string{"CreateInstanceTopicReqTopicOtherConfigs", string(data)}, " ")
}
