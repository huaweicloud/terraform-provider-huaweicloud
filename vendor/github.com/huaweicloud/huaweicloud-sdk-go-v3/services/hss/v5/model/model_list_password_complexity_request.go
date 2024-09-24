package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPasswordComplexityRequest Request Object
type ListPasswordComplexityRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器IP地址
	HostIp *string `json:"host_ip,omitempty"`

	// 主机id，不赋值时，查租户所有主机
	HostId *string `json:"host_id,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListPasswordComplexityRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPasswordComplexityRequest struct{}"
	}

	return strings.Join([]string{"ListPasswordComplexityRequest", string(data)}, " ")
}
