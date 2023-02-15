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

	// 防护状态，包含如下2种。 - closed ：未防护。 - opened ：防护中。
	ProtectStatus *string `json:"protect_status,omitempty"`

	// 服务器组ID
	GroupId *string `json:"group_id,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 策略ID
	ProtectPolicyId *string `json:"protect_policy_id,omitempty"`

	// 策略名称
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
}

func (o ProtectionServerInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtectionServerInfo struct{}"
	}

	return strings.Join([]string{"ProtectionServerInfo", string(data)}, " ")
}
