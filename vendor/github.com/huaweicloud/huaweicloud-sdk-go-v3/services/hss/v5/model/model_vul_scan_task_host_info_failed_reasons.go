package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VulScanTaskHostInfoFailedReasons struct {

	// 扫描失败的漏洞类型，包含如下：   -linux_vul : linux漏洞   -windows_vul : windows漏洞   -web_cms : Web-CMS漏洞   -app_vul : 应用漏洞   -urgent_vul : 应急漏洞
	VulType *string `json:"vul_type,omitempty"`

	// 扫描失败的原因
	FailedReason *string `json:"failed_reason,omitempty"`
}

func (o VulScanTaskHostInfoFailedReasons) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulScanTaskHostInfoFailedReasons struct{}"
	}

	return strings.Join([]string{"VulScanTaskHostInfoFailedReasons", string(data)}, " ")
}
