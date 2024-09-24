package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAutopilotReleaseResponse Response Object
type CreateAutopilotReleaseResponse struct {

	// 模板名称
	ChartName *string `json:"chart_name,omitempty"`

	// 是否公开模板
	ChartPublic *bool `json:"chart_public,omitempty"`

	// 模板版本
	ChartVersion *string `json:"chart_version,omitempty"`

	// 集群ID
	ClusterId *string `json:"cluster_id,omitempty"`

	// 集群名称
	ClusterName *string `json:"cluster_name,omitempty"`

	// 创建时间
	CreateAt *string `json:"create_at,omitempty"`

	// 模板实例描述
	Description *string `json:"description,omitempty"`

	// 模板实例名称
	Name *string `json:"name,omitempty"`

	// 模板实例所在的命名空间
	Namespace *string `json:"namespace,omitempty"`

	// 模板实例参数
	Parameters *string `json:"parameters,omitempty"`

	// 模板实例需要的资源
	Resources *string `json:"resources,omitempty"`

	// 模板实例状态
	Status *string `json:"status,omitempty"`

	// 模板实例状态描述
	StatusDescription *string `json:"status_description,omitempty"`

	// 更新时间
	UpdateAt *string `json:"update_at,omitempty"`

	// 模板实例的值
	Values *string `json:"values,omitempty"`

	// 模板实例版本
	Version        *int32 `json:"version,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateAutopilotReleaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAutopilotReleaseResponse struct{}"
	}

	return strings.Join([]string{"CreateAutopilotReleaseResponse", string(data)}, " ")
}
