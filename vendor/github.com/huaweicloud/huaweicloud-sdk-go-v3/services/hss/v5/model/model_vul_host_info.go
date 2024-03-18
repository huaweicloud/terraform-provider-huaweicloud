package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// VulHostInfo 软件漏洞列表
type VulHostInfo struct {

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 危险程度   - Critical : 漏洞cvss评分大于等于9；对应控制台页面的高危   - High : 漏洞cvss评分大于等于7，小于9；对应控制台页面的中危   - Medium : 漏洞cvss评分大于等于4，小于7；对应控制台页面的中危   - Low : 漏洞cvss评分小于4；对应控制台页面的低危
	SeverityLevel *string `json:"severity_level,omitempty"`

	// 受影响资产名称
	HostName *string `json:"host_name,omitempty"`

	// 受影响资产ip
	HostIp *string `json:"host_ip,omitempty"`

	// 主机对应的agent id
	AgentId *string `json:"agent_id,omitempty"`

	// 漏洞cve数
	CveNum *int32 `json:"cve_num,omitempty"`

	// cve列表
	CveIdList *[]string `json:"cve_id_list,omitempty"`

	// 漏洞状态   - vul_status_unfix : 未处理   - vul_status_ignored : 已忽略   - vul_status_verified : 验证中   - vul_status_fixing : 修复中   - vul_status_fixed : 修复成功   - vul_status_reboot : 修复成功待重启   - vul_status_failed : 修复失败   - vul_status_fix_after_reboot : 请重启主机再次修复
	Status *string `json:"status,omitempty"`

	// 修复命令行
	RepairCmd *string `json:"repair_cmd,omitempty"`

	// 应用软件的路径（只有应用漏洞有该字段）
	AppPath *string `json:"app_path,omitempty"`

	// 地域
	RegionName *string `json:"region_name,omitempty"`

	// 服务器公网ip
	PublicIp *string `json:"public_ip,omitempty"`

	// 服务器私网ip
	PrivateIp *string `json:"private_ip,omitempty"`

	// 服务器组id
	GroupId *string `json:"group_id,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 操作系统
	OsType *string `json:"os_type,omitempty"`

	// 资产重要性，包含如下3种   - important ：重要资产   - common ：一般资产   - test ：测试资产
	AssetValue *string `json:"asset_value,omitempty"`

	// 是否影响业务
	IsAffectBusiness *bool `json:"is_affect_business,omitempty"`

	// 首次扫描时间
	FirstScanTime *int64 `json:"first_scan_time,omitempty"`

	// 扫描时间
	ScanTime *int64 `json:"scan_time,omitempty"`

	// 是否可以回滚到修复漏洞时创建的备份
	SupportRestore *bool `json:"support_restore,omitempty"`
}

func (o VulHostInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulHostInfo struct{}"
	}

	return strings.Join([]string{"VulHostInfo", string(data)}, " ")
}
