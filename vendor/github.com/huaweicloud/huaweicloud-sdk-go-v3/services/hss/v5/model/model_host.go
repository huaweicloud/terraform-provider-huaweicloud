package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Host struct {

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器ID
	HostId *string `json:"host_id,omitempty"`

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// 私有IP地址
	PrivateIp *string `json:"private_ip,omitempty"`

	// 弹性公网IP地址
	PublicIp *string `json:"public_ip,omitempty"`

	// 企业项目ID
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 所属企业项目名称
	EnterpriseProjectName *string `json:"enterprise_project_name,omitempty"`

	// 服务器状态，包含如下4种。   - ACTIVE ：运行中。   - SHUTOFF ：关机。   - BUILDING ：创建中。   - ERROR ：故障。
	HostStatus *string `json:"host_status,omitempty"`

	// Agent状态，包含如下5种。   - installed ：已安装。   - not_installed ：未安装。   - online ：在线。   - offline ：离线。   - install_failed ：安装失败。   - installing ：安装中。
	AgentStatus *string `json:"agent_status,omitempty"`

	// 安装结果，包含如下12种。   - install_succeed ：安装成功。   - network_access_timeout ：网络不通，访问超时。   - invalid_port ：无效端口。   - auth_failed ：认证错误，口令不正确。   - permission_denied ：权限错误，被拒绝。   - no_available_vpc ：没有相同VPC的agent在线虚拟机。   - install_exception ：安装异常。   - invalid_param ：参数错误。   - install_failed ：安装失败。   - package_unavailable ：安装包失效。   - os_type_not_support ：系统类型错误。   - os_arch_not_support ：架构类型错误。
	InstallResultCode *string `json:"install_result_code,omitempty"`

	// 主机开通的版本，包含如下7种输入。   - hss.version.null ：无。   - hss.version.basic ：基础版。   - hss.version.advanced ：专业版。   - hss.version.enterprise ：企业版。   - hss.version.premium ：旗舰版。   - hss.version.wtp ：网页防篡改版。   - hss.version.container.enterprise ：容器版。
	Version *string `json:"version,omitempty"`

	// 防护状态，包含如下2种。 - closed ：未防护。 - opened ：防护中。
	ProtectStatus *string `json:"protect_status,omitempty"`

	// 系统镜像
	OsImage *string `json:"os_image,omitempty"`

	// 操作系统类型，包含如下2种。   - Linux ：Linux。   - Windows ：Windows。
	OsType *string `json:"os_type,omitempty"`

	// 操作系统位数
	OsBit *string `json:"os_bit,omitempty"`

	// 云主机安全检测结果，包含如下4种。 - undetected ：未检测。 - clean ：无风险。 - risk ：有风险。 - scanning ：检测中。
	DetectResult *string `json:"detect_result,omitempty"`

	// 试用版到期时间（-1表示非试用版配额，当值不为-1时为试用版本过期时间）
	ExpireTime *int64 `json:"expire_time,omitempty"`

	// 收费模式，包含如下2种。   - packet_cycle ：包年/包月。   - on_demand ：按需。
	ChargingMode *string `json:"charging_mode,omitempty"`

	// 主机安全配额ID（UUID）
	ResourceId *string `json:"resource_id,omitempty"`

	// 是否非华为云机器
	OutsideHost *bool `json:"outside_host,omitempty"`

	// 服务器组ID
	GroupId *string `json:"group_id,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 策略组ID
	PolicyGroupId *string `json:"policy_group_id,omitempty"`

	// 策略组名称
	PolicyGroupName *string `json:"policy_group_name,omitempty"`

	// 资产风险
	Asset *int32 `json:"asset,omitempty"`

	// 漏洞风险总数，包含Linux软件漏洞、Windows系统漏洞、Web-CMS漏洞、应用漏洞
	Vulnerability *int32 `json:"vulnerability,omitempty"`

	// 基线风险总数，包含配置风险、弱口令
	Baseline *int32 `json:"baseline,omitempty"`

	// 入侵风险总数
	Intrusion *int32 `json:"intrusion,omitempty"`

	// 资产重要性，包含如下4种   - important ：重要资产   - common ：一般资产   - test ：测试资产
	AssetValue *string `json:"asset_value,omitempty"`

	// 标签列表
	Labels *[]string `json:"labels,omitempty"`

	// agent安装时间，采用时间戳，默认毫秒，
	AgentCreateTime *int64 `json:"agent_create_time,omitempty"`

	// agent状态修改时间，采用时间戳，默认毫秒，
	AgentUpdateTime *int64 `json:"agent_update_time,omitempty"`

	// agent版本
	AgentVersion *string `json:"agent_version,omitempty"`

	// 升级状态，包含如下4种。   - not_upgrade ：未升级，也就是默认状态，客户还没有给这台机器下发过升级。   - upgrading ：正在升级中。   - upgrade_failed ：升级失败。   - upgrade_succeed ：升级成功。
	UpgradeStatus *string `json:"upgrade_status,omitempty"`

	// 升级失败原因，只有当 upgrade_status 为 upgrade_failed 时才显示，包含如下6种。   - package_unavailable ：升级包解析失败，升级文件有错误。   - network_access_timeout ：下载升级包失败，网络异常。   - agent_offline ：agent离线。   - hostguard_abnormal ：agent工作进程异常。   - insufficient_disk_space ：磁盘空间不足。   - failed_to_replace_file ：替换文件失败。
	UpgradeResultCode *string `json:"upgrade_result_code,omitempty"`

	// 该服务器agent是否可升级
	Upgradable *bool `json:"upgradable,omitempty"`

	// 开启防护时间，采用时间戳，默认毫秒，
	OpenTime *int64 `json:"open_time,omitempty"`

	// 防护是否中断
	ProtectInterrupt *bool `json:"protect_interrupt,omitempty"`
}

func (o Host) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Host struct{}"
	}

	return strings.Join([]string{"Host", string(data)}, " ")
}
