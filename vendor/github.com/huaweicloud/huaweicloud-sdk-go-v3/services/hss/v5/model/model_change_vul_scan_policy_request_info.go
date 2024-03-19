package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ChangeVulScanPolicyRequestInfo struct {

	// 扫描周期 - one_day : 每天 - three_day : 每三天 - one_week : 每周
	ScanPeriod string `json:"scan_period"`

	// 扫描主机的范围，包含如下：   -all_host : 扫描全部主机   -specific_host : 扫描指定主机
	ScanRangeType string `json:"scan_range_type"`

	// 主机ID列表；当scan_range_type的值为specific_host时必填
	HostIds *[]string `json:"host_ids,omitempty"`

	// 扫描的漏洞类型列表
	ScanVulTypes *[]string `json:"scan_vul_types,omitempty"`

	// 扫描策略状态，包含如下：   -open : 开启   -close : 关闭
	Status string `json:"status"`
}

func (o ChangeVulScanPolicyRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeVulScanPolicyRequestInfo struct{}"
	}

	return strings.Join([]string{"ChangeVulScanPolicyRequestInfo", string(data)}, " ")
}
