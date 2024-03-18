package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type WorkFlowSpec struct {

	// 集群ID，资源唯一标识，创建成功后自动生成，填写无效
	ClusterID *string `json:"clusterID,omitempty"`

	// 本次集群升级的当前版本
	ClusterVersion *string `json:"clusterVersion,omitempty"`

	// 本次集群升级的目标版本
	TargetVersion string `json:"targetVersion"`
}

func (o WorkFlowSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "WorkFlowSpec struct{}"
	}

	return strings.Join([]string{"WorkFlowSpec", string(data)}, " ")
}
