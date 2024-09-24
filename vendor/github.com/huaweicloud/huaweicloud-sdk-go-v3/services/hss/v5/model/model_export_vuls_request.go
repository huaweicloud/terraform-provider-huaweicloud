package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExportVulsRequest Request Object
type ExportVulsRequest struct {

	// 企业租户ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 漏洞类型，包含如下：   -linux_vul : linux漏洞   -windows_vul : windows漏洞   -web_cms : Web-CMS漏洞   -app_vul : 应用漏洞   -urgent_vul : 应急漏洞
	Type *string `json:"type,omitempty"`

	// 漏洞ID
	VulId *string `json:"vul_id,omitempty"`

	// 漏洞名称
	VulName *string `json:"vul_name,omitempty"`

	// 主机id，导出单台主机漏洞时会用到
	HostId *string `json:"host_id,omitempty"`

	// limit
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置，必须为数字，取值范围为大于或等于0，默认0
	Offset *int32 `json:"offset,omitempty"`

	// 修复优先级 Critical 紧急 High  高 Medium 中 Low 低
	RepairPriority *string `json:"repair_priority,omitempty"`

	// 处置状态，包含如下:   - unhandled ：未处理   - handled : 已处理
	HandleStatus *string `json:"handle_status,omitempty"`

	// 漏洞编号
	CveId *string `json:"cve_id,omitempty"`

	// 漏洞标签
	LabelList *string `json:"label_list,omitempty"`

	// 漏洞状态
	Status *string `json:"status,omitempty"`

	// 资产重要性 important common test
	AssetValue *string `json:"asset_value,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 导出数据条数
	ExportSize int32 `json:"export_size"`

	// 导出漏洞数据类别:   - vul ：漏洞   - host: 主机漏洞
	Category string `json:"category"`

	Body *ExportVulRequestBody `json:"body,omitempty"`
}

func (o ExportVulsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExportVulsRequest struct{}"
	}

	return strings.Join([]string{"ExportVulsRequest", string(data)}, " ")
}
