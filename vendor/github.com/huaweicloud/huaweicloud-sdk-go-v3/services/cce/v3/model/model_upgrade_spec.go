package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpgradeSpec struct {
	ClusterUpgradeAction *ClusterUpgradeAction `json:"clusterUpgradeAction,omitempty"`
}

func (o UpgradeSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeSpec struct{}"
	}

	return strings.Join([]string{"UpgradeSpec", string(data)}, " ")
}
