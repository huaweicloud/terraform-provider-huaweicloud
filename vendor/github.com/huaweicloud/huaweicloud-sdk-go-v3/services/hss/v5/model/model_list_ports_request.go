package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPortsRequest Request Object
type ListPortsRequest struct {

	// 主机id
	HostId string `json:"host_id"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 端口号
	Port *int32 `json:"port,omitempty"`

	// 端口类型：目前包括TCP，UDP两种
	Type *string `json:"type,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 类别，默认为host，包含如下： - host：主机 - container：容器
	Category *string `json:"category,omitempty"`
}

func (o ListPortsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPortsRequest struct{}"
	}

	return strings.Join([]string{"ListPortsRequest", string(data)}, " ")
}
