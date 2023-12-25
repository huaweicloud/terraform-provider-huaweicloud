package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ClusterConfigurationsSpecPackages struct {

	// 组件名称
	Name *string `json:"name,omitempty"`

	// 组件配置项详情
	Configurations *[]ConfigurationItem `json:"configurations,omitempty"`
}

func (o ClusterConfigurationsSpecPackages) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterConfigurationsSpecPackages struct{}"
	}

	return strings.Join([]string{"ClusterConfigurationsSpecPackages", string(data)}, " ")
}
