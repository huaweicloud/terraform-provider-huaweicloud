package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HostVulInfo struct {

	// 漏洞名称
	VulName *string `json:"vul_name,omitempty"`

	// 漏洞ID
	VulId *string `json:"vul_id,omitempty"`

	// 漏洞标签列表
	LabelList *[]string `json:"label_list,omitempty"`

	// 修复紧急度，包括如下：   - immediate_repair : 尽快修复   - delay_repair : 延后修复   - not_needed_repair : 暂可不修复
	RepairNecessity *string `json:"repair_necessity,omitempty"`

	// 最近扫描时间
	ScanTime *int64 `json:"scan_time,omitempty"`

	// 漏洞类型，包含如下：   -linux_vul : linux漏洞   -windows_vul : windows漏洞   -web_cms : Web-CMS漏洞   -app_vul : 应用漏洞
	Type *string `json:"type,omitempty"`

	// 服务器上受该漏洞影响的软件列表
	AppList *[]HostVulInfoAppList `json:"app_list,omitempty"`

	// 危险程度   - Critical : 致命   - High : 高危   - Medium : 中危   - Low : 低危
	SeverityLevel *string `json:"severity_level,omitempty"`

	// 解决方案
	SolutionDetail *string `json:"solution_detail,omitempty"`

	// URL链接
	Url *string `json:"url,omitempty"`

	// 漏洞描述
	Description *string `json:"description,omitempty"`

	// 修复命令行
	RepairCmd *string `json:"repair_cmd,omitempty"`

	// 漏洞状态   - vul_status_unfix : 未处理   - vul_status_ignored : 已忽略   - vul_status_verified : 验证中   - vul_status_fixing : 修复中   - vul_status_fixed : 修复成功   - vul_status_reboot : 修复成功待重启   - vul_status_failed : 修复失败   - vul_status_fix_after_reboot : 请重启主机再次修复
	Status *string `json:"status,omitempty"`

	// HSS全网修复该漏洞的次数
	RepairSuccessNum *int32 `json:"repair_success_num,omitempty"`

	// CVE列表
	CveList *[]HostVulInfoCveList `json:"cve_list,omitempty"`

	// 是否影响业务
	IsAffectBusiness *bool `json:"is_affect_business,omitempty"`

	// 首次扫描时间
	FirstScanTime *int64 `json:"first_scan_time,omitempty"`
}

func (o HostVulInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostVulInfo struct{}"
	}

	return strings.Join([]string{"HostVulInfo", string(data)}, " ")
}
