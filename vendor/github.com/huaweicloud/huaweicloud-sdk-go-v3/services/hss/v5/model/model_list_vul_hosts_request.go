package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListVulHostsRequest struct {

	// 企业租户ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 漏洞ID
	VulId string `json:"vul_id"`

	// 漏洞类型   - linux_vul : 漏洞类型-linux漏洞   - windows_vul : 漏洞类型-windows漏洞
	Type string `json:"type"`

	// 受影响资产名称
	HostName *string `json:"host_name,omitempty"`

	// 受影响资产ip
	HostIp *string `json:"host_ip,omitempty"`

	// 漏洞状态   - vul_status_unfix : 未处理   - vul_status_ignored : 已忽略   - vul_status_verified : 验证中   - vul_status_fixing : 修复中   - vul_status_fixed : 修复成功   - vul_status_reboot : 修复成功待重启   - vul_status_failed : 修复失败   - vul_status_fix_after_reboot : 请重启主机再次修复
	Status *string `json:"status,omitempty"`

	// 每页条数
	Limit *int32 `json:"limit,omitempty"`

	// 偏移
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListVulHostsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVulHostsRequest struct{}"
	}

	return strings.Join([]string{"ListVulHostsRequest", string(data)}, " ")
}
