package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradePath 升级路径
type UpgradePath struct {

	// 集群版本，v1.19及以下集群形如v1.19.16-r20，v1.21及以上形如v1.21,v1.23，详细请参考CCE集群版本号说明。
	Version *string `json:"version,omitempty"`

	// CCE集群平台版本号，表示集群版本(version)下的内部版本。用于跟踪某一集群版本内的迭代，集群版本内唯一，跨集群版本重新计数。   platformVersion格式为：cce.X.Y   - X: 表示内部特性版本。集群版本中特性或者补丁修复，或者OS支持等变更场景。其值从1开始单调递增。  - Y: 表示内部特性版本的补丁版本。仅用于特性版本上线后的软件包更新，不涉及其他修改。其值从0开始单调递增。
	PlatformVersion *string `json:"platformVersion,omitempty"`

	// 可升级的目标版本集合
	TargetVersions *[]string `json:"targetVersions,omitempty"`
}

func (o UpgradePath) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradePath struct{}"
	}

	return strings.Join([]string{"UpgradePath", string(data)}, " ")
}
