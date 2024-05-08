package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// VulInfo 漏洞列表
type VulInfo struct {

	// 漏洞名称
	VulName *string `json:"vul_name,omitempty"`

	// 漏洞ID
	VulId *string `json:"vul_id,omitempty"`

	// 漏洞标签
	LabelList *[]string `json:"label_list,omitempty"`

	// 修复必要性   - Critical : 漏洞cvss评分大于等于9；对应控制台页面的高危   - High : 漏洞cvss评分大于等于7，小于9；对应控制台页面的中危   - Medium : 漏洞cvss评分大于等于4，小于7；对应控制台页面的中危   - Low : 漏洞cvss评分小于4；对应控制台页面的低危
	RepairNecessity *string `json:"repair_necessity,omitempty"`

	// 漏洞级别   - Critical : 漏洞cvss评分大于等于9；对应控制台页面的高危   - High : 漏洞cvss评分大于等于7，小于9；对应控制台页面的中危   - Medium : 漏洞cvss评分大于等于4，小于7；对应控制台页面的中危   - Low : 漏洞cvss评分小于4；对应控制台页面的低危
	SeverityLevel *string `json:"severity_level,omitempty"`

	// 受影响服务器台数
	HostNum *int32 `json:"host_num,omitempty"`

	// 未处理主机台数，除已忽略和已修复的主机数量
	UnhandleHostNum *int32 `json:"unhandle_host_num,omitempty"`

	// 最近扫描时间，时间戳单位：毫秒
	ScanTime *int64 `json:"scan_time,omitempty"`

	// 修复漏洞的指导意见
	SolutionDetail *string `json:"solution_detail,omitempty"`

	// URL链接
	Url *string `json:"url,omitempty"`

	// 漏洞描述
	Description *string `json:"description,omitempty"`

	// 漏洞类型，包含如下：   -linux_vul : linux漏洞   -windows_vul : windows漏洞   -web_cms : Web-CMS漏洞   -app_vul : 应用漏洞
	Type *string `json:"type,omitempty"`

	// 可处置该漏洞的主机列表
	HostIdList *[]string `json:"host_id_list,omitempty"`

	// CVE列表
	CveList *[]VulInfoCveList `json:"cve_list,omitempty"`

	// 补丁地址
	PatchUrl *string `json:"patch_url,omitempty"`

	// 修复优先级 Critical 紧急 High 高 Medium 中 Low 低
	RepairPriority *string `json:"repair_priority,omitempty"`

	HostsNum *VulnerabilityHostNumberInfo `json:"hosts_num,omitempty"`

	// 修复成功次数
	RepairSuccessNum *int32 `json:"repair_success_num,omitempty"`

	// 修复数量
	FixedNum *int64 `json:"fixed_num,omitempty"`

	// 忽略数量
	IgnoredNum *int64 `json:"ignored_num,omitempty"`

	// 验证数量
	VerifyNum *int32 `json:"verify_num,omitempty"`

	// 修复优先级，每个修复优先级对应的主机数量
	RepairPriorityList *[]RepairPriorityListInfo `json:"repair_priority_list,omitempty"`
}

func (o VulInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulInfo struct{}"
	}

	return strings.Join([]string{"VulInfo", string(data)}, " ")
}
