package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ContainerBaseInfo 容器基本信息
type ContainerBaseInfo struct {

	// 容器ID
	ContainerId *string `json:"container_id,omitempty"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 容器状态，包含以下几种： -Running 运行中 -Terminated 终止 -Waiting 等待
	Status *string `json:"status,omitempty"`

	// 创建时间
	CreateTime *int64 `json:"create_time,omitempty"`

	// cpu限制
	CpuLimit *string `json:"cpu_limit,omitempty"`

	// 内存限制
	MemoryLimit *string `json:"memory_limit,omitempty"`

	// 重启次数
	RestartCount *int32 `json:"restart_count,omitempty"`

	// 所属pod名称
	PodName *string `json:"pod_name,omitempty"`

	// 所属集群
	ClusterName *string `json:"cluster_name,omitempty"`

	// 集群id
	ClusterId *string `json:"cluster_id,omitempty"`

	// 集群类型，包含以下几种： -k8s 原生集群 -cce CCE集群 -ali 阿里云集群 -tencent 腾讯云集群 -azure 微软云集群 -aws 亚马逊集群 -self_built_hw 华为云自建集群 -self_built_idc IDC自建集群
	ClusterType *string `json:"cluster_type,omitempty"`

	// 是否有风险
	Risky *bool `json:"risky,omitempty"`

	// 低危风险数量
	LowRisk *int32 `json:"low_risk,omitempty"`

	// 中危风险数量
	MediumRisk *int32 `json:"medium_risk,omitempty"`

	// 高危风险数量
	HighRisk *int32 `json:"high_risk,omitempty"`

	// 致命风险数量
	FatalRisk *int32 `json:"fatal_risk,omitempty"`
}

func (o ContainerBaseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ContainerBaseInfo struct{}"
	}

	return strings.Join([]string{"ContainerBaseInfo", string(data)}, " ")
}
