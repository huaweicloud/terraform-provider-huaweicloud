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

	// 修复必要性
	RepairNecessity *string `json:"repair_necessity,omitempty"`

	// 漏洞级别
	SeverityLevel *string `json:"severity_level,omitempty"`

	// 受影响服务器台数
	HostNum *int32 `json:"host_num,omitempty"`

	// 未处理服务器台数
	UnhandleHostNum *int32 `json:"unhandle_host_num,omitempty"`

	// 最近扫描时间
	ScanTime *int64 `json:"scan_time,omitempty"`

	// 解决方案
	SolutionDetail *string `json:"solution_detail,omitempty"`

	// URL链接
	Url *string `json:"url,omitempty"`

	// 漏洞描述
	Description *string `json:"description,omitempty"`

	// 漏洞类型，包含如下：   -linux_vul : linux漏洞   -windows_vul : windows漏洞   -web_cms : Web-CMS漏洞   -app_vul : 应用漏洞
	Type *string `json:"type,omitempty"`

	// 主机列表
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
}

func (o VulInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulInfo struct{}"
	}

	return strings.Join([]string{"VulInfo", string(data)}, " ")
}
