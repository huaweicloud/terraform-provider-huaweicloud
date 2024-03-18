package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutoLaunchsRequest Request Object
type ListAutoLaunchsRequest struct {

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 自启动项名称
	Name *string `json:"name,omitempty"`

	// 主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 自启动项类型
	Type *string `json:"type,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`

	// 是否模糊匹配，默认false表示精确匹配
	PartMatch *bool `json:"part_match,omitempty"`
}

func (o ListAutoLaunchsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutoLaunchsRequest struct{}"
	}

	return strings.Join([]string{"ListAutoLaunchsRequest", string(data)}, " ")
}
