package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PostcheckSpec struct {

	// 集群ID
	ClusterID *string `json:"clusterID,omitempty"`

	// 集群升级源版本
	ClusterVersion *string `json:"clusterVersion,omitempty"`

	// 集群升级目标版本
	TargetVersion *string `json:"targetVersion,omitempty"`
}

func (o PostcheckSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostcheckSpec struct{}"
	}

	return strings.Join([]string{"PostcheckSpec", string(data)}, " ")
}
