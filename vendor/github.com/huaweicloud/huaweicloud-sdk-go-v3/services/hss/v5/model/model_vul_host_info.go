package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 软件漏洞列表
type VulHostInfo struct {

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 危险程度   - Critical : 高危   - High : 中危   - Medium : 中危   - Low : 低危
	SeverityLevel *string `json:"severity_level,omitempty"`

	// 受影响资产名称
	HostName *string `json:"host_name,omitempty"`

	// 受影响资产ip
	HostIp *string `json:"host_ip,omitempty"`

	// 漏洞cve数
	CveNum *int32 `json:"cve_num,omitempty"`

	// cve列表
	CveIdList *[]string `json:"cve_id_list,omitempty"`

	// 漏洞状态   - vul_status_unfix : 未处理   - vul_status_ignored : 已忽略   - vul_status_verified : 验证中   - vul_status_fixing : 修复中   - vul_status_fixed : 修复成功   - vul_status_reboot : 修复成功待重启   - vul_status_failed : 修复失败   - vul_status_fix_after_reboot : 请重启主机再次修复
	Status *string `json:"status,omitempty"`

	// 修复命令行
	RepairCmd *string `json:"repair_cmd,omitempty"`
}

func (o VulHostInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulHostInfo struct{}"
	}

	return strings.Join([]string{"VulHostInfo", string(data)}, " ")
}
