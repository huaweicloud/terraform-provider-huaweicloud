package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRiskConfigHostsRequest Request Object
type ListRiskConfigHostsRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 基线名称
	CheckName string `json:"check_name"`

	// 标准类型，包含如下: - cn_standard : 等保合规标准 - hw_standard : 华为标准 - qt_standard : 青腾标准
	Standard string `json:"standard"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器IP地址
	HostIp *string `json:"host_ip,omitempty"`

	// 每页数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置，必须为数字，取值范围为大于或等于0，默认0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListRiskConfigHostsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRiskConfigHostsRequest struct{}"
	}

	return strings.Join([]string{"ListRiskConfigHostsRequest", string(data)}, " ")
}
