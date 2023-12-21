package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ClusterUpgradeResponseAction struct {

	// 当前集群版本
	Version *string `json:"version,omitempty"`

	// 目标集群版本，例如\"v1.23\"
	TargetVersion *string `json:"targetVersion,omitempty"`

	// 目标集群的平台版本号，表示集群版本(version)下的内部版本，不支持用户指定。
	TargetPlatformVersion *string `json:"targetPlatformVersion,omitempty"`

	Strategy *UpgradeStrategy `json:"strategy,omitempty"`

	// 升级过程中指定的集群配置
	Config *interface{} `json:"config,omitempty"`
}

func (o ClusterUpgradeResponseAction) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterUpgradeResponseAction struct{}"
	}

	return strings.Join([]string{"ClusterUpgradeResponseAction", string(data)}, " ")
}
