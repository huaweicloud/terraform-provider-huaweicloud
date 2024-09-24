package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AccountResponseInfo 事件列表详情
type AccountResponseInfo struct {

	// 账号名称
	AccountName *string `json:"account_name,omitempty"`

	// 账号Id
	AccountId *string `json:"account_id,omitempty"`

	// 组织Id
	OrganizationId *string `json:"organization_id,omitempty"`

	// 项目Id
	ProjectId *string `json:"project_id,omitempty"`

	// 项目名称
	ProjectName *string `json:"project_name,omitempty"`

	// 主机数量
	HostNum *int32 `json:"host_num,omitempty"`

	// 漏洞风险数量
	VulnerabilityNum *int32 `json:"vulnerability_num,omitempty"`

	// 基线检测风险数量
	BaselineNum *int32 `json:"baseline_num,omitempty"`

	// 安全告警风险数量
	IntrusionNum *int32 `json:"intrusion_num,omitempty"`
}

func (o AccountResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AccountResponseInfo struct{}"
	}

	return strings.Join([]string{"AccountResponseInfo", string(data)}, " ")
}
