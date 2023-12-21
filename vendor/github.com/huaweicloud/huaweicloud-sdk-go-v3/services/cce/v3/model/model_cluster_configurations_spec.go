package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClusterConfigurationsSpec Configuration的规格信息
type ClusterConfigurationsSpec struct {

	// 组件配置项列表
	Packages []ClusterConfigurationsSpecPackages `json:"packages"`
}

func (o ClusterConfigurationsSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterConfigurationsSpec struct{}"
	}

	return strings.Join([]string{"ClusterConfigurationsSpec", string(data)}, " ")
}
