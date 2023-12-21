package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeTaskSpec 升级任务属性
type UpgradeTaskSpec struct {

	// 升级前集群版本
	Version *string `json:"version,omitempty"`

	// 升级的目标集群版本
	TargetVersion *string `json:"targetVersion,omitempty"`

	// 升级任务附属信息
	Items *interface{} `json:"items,omitempty"`
}

func (o UpgradeTaskSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeTaskSpec struct{}"
	}

	return strings.Join([]string{"UpgradeTaskSpec", string(data)}, " ")
}
