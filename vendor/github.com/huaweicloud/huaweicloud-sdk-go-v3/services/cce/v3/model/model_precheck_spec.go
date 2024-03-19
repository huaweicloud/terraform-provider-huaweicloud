package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PrecheckSpec struct {

	// 集群ID
	ClusterID *string `json:"clusterID,omitempty"`

	// 集群版本
	ClusterVersion *string `json:"clusterVersion,omitempty"`

	// 升级目标版本
	TargetVersion *string `json:"targetVersion,omitempty"`

	// 跳过检查的项目列表
	SkippedCheckItemList *[]SkippedCheckItemList `json:"skippedCheckItemList,omitempty"`
}

func (o PrecheckSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrecheckSpec struct{}"
	}

	return strings.Join([]string{"PrecheckSpec", string(data)}, " ")
}
