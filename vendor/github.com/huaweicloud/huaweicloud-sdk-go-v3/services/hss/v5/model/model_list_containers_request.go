package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListContainersRequest Request Object
type ListContainersRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`

	// 所属Pod名称
	PodName *string `json:"pod_name,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 是否是集群纳管的容器
	ClusterContainer *bool `json:"cluster_container,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListContainersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListContainersRequest struct{}"
	}

	return strings.Join([]string{"ListContainersRequest", string(data)}, " ")
}
