package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRiskConfigHostsRequest Request Object
type ListRiskConfigHostsRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 配置检查（基线）的名称，例如SSH、CentOS 7、Windows
	CheckName string `json:"check_name"`

	// 标准类型，包含如下: - cn_standard : 等保合规标准 - hw_standard : 云安全实践标准
	Standard string `json:"standard"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器IP地址
	HostIp *string `json:"host_ip,omitempty"`

	// 每页数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListRiskConfigHostsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRiskConfigHostsRequest struct{}"
	}

	return strings.Join([]string{"ListRiskConfigHostsRequest", string(data)}, " ")
}
