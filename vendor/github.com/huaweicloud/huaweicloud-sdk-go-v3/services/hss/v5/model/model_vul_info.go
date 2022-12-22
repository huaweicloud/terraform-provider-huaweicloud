package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 漏洞列表
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

	// 漏洞类型，包含如下：   -linux_vul : linux漏洞   -windows_vul : windows漏洞   -web_cms : Web-CMS漏洞
	Type *string `json:"type,omitempty"`

	// 主机列表
	HostIdList *[]string `json:"host_id_list,omitempty"`
}

func (o VulInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulInfo struct{}"
	}

	return strings.Join([]string{"VulInfo", string(data)}, " ")
}
