package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRiskConfigsRequest Request Object
type ListRiskConfigsRequest struct {

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 配置检查（基线）的名称，例如SSH、CentOS 7、Windows
	CheckName *string `json:"check_name,omitempty"`

	// 策略组ID
	GroupId *string `json:"group_id,omitempty"`

	// 风险等级，包含如下:   - Security : 安全   - Low : 低危   - Medium : 中危   - High : 高危
	Severity *string `json:"severity,omitempty"`

	// 标准类型，包含如下:   - cn_standard : 等保合规标准   - hw_standard : 云安全实践标准
	Standard *string `json:"standard,omitempty"`

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListRiskConfigsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRiskConfigsRequest struct{}"
	}

	return strings.Join([]string{"ListRiskConfigsRequest", string(data)}, " ")
}
