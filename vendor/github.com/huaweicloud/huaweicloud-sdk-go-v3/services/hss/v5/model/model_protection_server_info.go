package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProtectionServerInfo struct {

	// 服务器ID
	HostId *string `json:"host_id,omitempty"`

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 弹性公网IP地址
	HostIp *string `json:"host_ip,omitempty"`

	// 私有IP地址
	PrivateIp *string `json:"private_ip,omitempty"`

	// 操作系统类型，包含如下2种。   - Linux ：Linux。   - Windows ：Windows。
	OsType *string `json:"os_type,omitempty"`

	// 系统名称
	OsName *string `json:"os_name,omitempty"`

	// 服务器状态，包含如下2种。   - ACTIVE ：运行中。   - SHUTOFF ：关机。
	HostStatus *string `json:"host_status,omitempty"`

	// 勒索防护状态，包含如下4种。   - closed ：关闭。   - opened ：开启。   - opening ：开启中。   - closing ：关闭中。
	RansomProtectionStatus *string `json:"ransom_protection_status,omitempty"`

	// agent版本
	AgentVersion *string `json:"agent_version,omitempty"`

	// 防护状态，包含如下2种。 - closed ：未防护。 - opened ：防护中。
	ProtectStatus *string `json:"protect_status,omitempty"`

	// 服务器组ID
	GroupId *string `json:"group_id,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 防护策略ID
	ProtectPolicyId *string `json:"protect_policy_id,omitempty"`

	// 防护策略名称
	ProtectPolicyName *string `json:"protect_policy_name,omitempty"`

	BackupError *ProtectionServerInfoBackupError `json:"backup_error,omitempty"`

	// 是否开启备份，包含如下3种。   - failed_to_turn_on_backup: 无法开启备份   - closed ：关闭。   - opened ：开启。
	BackupProtectionStatus *string `json:"backup_protection_status,omitempty"`

	// 防护事件数
	CountProtectEvent *int32 `json:"count_protect_event,omitempty"`

	// 已有备份数
	CountBackuped *int32 `json:"count_backuped,omitempty"`

	// Agent状态
	AgentStatus *string `json:"agent_status,omitempty"`

	// 主机开通的版本，包含如下7种输入。   - hss.version.null ：无。   - hss.version.basic ：基础版。   - hss.version.advanced ：专业版。   - hss.version.enterprise ：企业版。   - hss.version.premium ：旗舰版。   - hss.version.wtp ：网页防篡改版。   - hss.version.container.enterprise ：容器版。
	Version *string `json:"version,omitempty"`

	// 服务器类型，包含如下3种输入。   - ecs ：ecs。   - outside ：线下主机。   - workspace ：云桌面。
	HostSource *string `json:"host_source,omitempty"`

	// 存储库ID
	VaultId *string `json:"vault_id,omitempty"`

	// 存储库名称
	VaultName *string `json:"vault_name,omitempty"`

	// 总容量，单位GB
	VaultSize *int32 `json:"vault_size,omitempty"`

	// 已使用容量，单位MB
	VaultUsed *int32 `json:"vault_used,omitempty"`

	// 已分配容量，单位GB，指绑定的服务器大小
	VaultAllocated *int32 `json:"vault_allocated,omitempty"`

	// 存储库创建模式，按需：post_paid，包周期：pre_paid
	VaultChargingMode *string `json:"vault_charging_mode,omitempty"`

	// 存储库状态。   - available ：可用。   - lock ：被锁定。   - frozen：冻结。   - deleting：删除中。   - error：错误。
	VaultStatus *string `json:"vault_status,omitempty"`

	// 备份策略ID，若为空，则为未绑定状态，若不为空，通过backup_policy_enabled字段判断策略是否启用
	BackupPolicyId *string `json:"backup_policy_id,omitempty"`

	// 备份策略名称
	BackupPolicyName *string `json:"backup_policy_name,omitempty"`

	// 策略是否启用
	BackupPolicyEnabled *bool `json:"backup_policy_enabled,omitempty"`

	// 已绑定服务器（个）
	ResourcesNum *int32 `json:"resources_num,omitempty"`
}

func (o ProtectionServerInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtectionServerInfo struct{}"
	}

	return strings.Join([]string{"ProtectionServerInfo", string(data)}, " ")
}
