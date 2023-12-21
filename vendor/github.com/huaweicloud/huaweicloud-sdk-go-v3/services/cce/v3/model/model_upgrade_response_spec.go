package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeResponseSpec 升级任务元数据
type UpgradeResponseSpec struct {
	ClusterUpgradeAction *ClusterUpgradeResponseAction `json:"clusterUpgradeAction,omitempty"`
}

func (o UpgradeResponseSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeResponseSpec struct{}"
	}

	return strings.Join([]string{"UpgradeResponseSpec", string(data)}, " ")
}
