package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListVulHostsRequest Request Object
type ListVulHostsRequest struct {

	// 企业项目ID，“0”表示默认企业项目，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 漏洞ID
	VulId string `json:"vul_id"`

	// 漏洞类型   - linux_vul : 漏洞类型-linux漏洞   - windows_vul : 漏洞类型-windows漏洞   - web_cms : Web-CMS漏洞   - app_vul : 应用漏洞   - urgent_vul : 应急漏洞
	Type string `json:"type"`

	// 受影响主机名称
	HostName *string `json:"host_name,omitempty"`

	// 受影响主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 漏洞状态   - vul_status_unfix : 未处理   - vul_status_ignored : 已忽略   - vul_status_verified : 验证中   - vul_status_fixing : 修复中   - vul_status_fixed : 修复成功   - vul_status_reboot : 修复成功待重启   - vul_status_failed : 修复失败   - vul_status_fix_after_reboot : 请重启主机再次修复
	Status *string `json:"status,omitempty"`

	// 每页条数
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 资产重要性 important:重要 common：一般 test：测试
	AssetValue *string `json:"asset_value,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 处置状态，包含如下:   - unhandled ：未处理   - handled : 已处理
	HandleStatus *string `json:"handle_status,omitempty"`

	// 危险程度 ，Critical，High，Medium，Low
	SeverityLevel *string `json:"severity_level,omitempty"`

	// 是否影响业务
	IsAffectBusiness *bool `json:"is_affect_business,omitempty"`

	// 修复优先级,包含如下 - Critical 紧急  - High 高 - Medium 中 - Low 低
	RepairPriority *string `json:"repair_priority,omitempty"`
}

func (o ListVulHostsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVulHostsRequest struct{}"
	}

	return strings.Join([]string{"ListVulHostsRequest", string(data)}, " ")
}
