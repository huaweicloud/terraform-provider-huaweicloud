package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AutopilotPackageConfiguration struct {

	// 组件名称
	Name *string `json:"name,omitempty"`

	// 组件配置项
	Configurations *[]AutopilotConfigurationItem `json:"configurations,omitempty"`
}

func (o AutopilotPackageConfiguration) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotPackageConfiguration struct{}"
	}

	return strings.Join([]string{"AutopilotPackageConfiguration", string(data)}, " ")
}
