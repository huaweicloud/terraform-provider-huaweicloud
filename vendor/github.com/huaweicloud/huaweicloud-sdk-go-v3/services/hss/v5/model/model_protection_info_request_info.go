package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProtectionInfoRequestInfo struct {

	// 操作系统，包含如下：   - Windows : 无需处理   - Linux : 已忽略
	OperatingSystem *string `json:"operating_system,omitempty"`

	// 勒索防护是否开启，包含如下：   - closed ：关闭。   - opened ：开启。
	RansomProtectionStatus *string `json:"ransom_protection_status,omitempty"`

	// 防护策略ID
	ProtectionPolicyId *string `json:"protection_policy_id,omitempty"`

	CreateProtectionPolicy *ProtectionProxyInfoRequestInfo `json:"create_protection_policy,omitempty"`

	// 是否服务器备份，包含如下：   - closed ：关闭。   - opened ：开启。
	BackupProtectionStatus *string `json:"backup_protection_status,omitempty"`

	// 备份策略ID
	BackupPolicyId *string `json:"backup_policy_id,omitempty"`

	BackupCycle *UpdateBackupPolicyRequestInfo `json:"backup_cycle,omitempty"`

	// 开启防护的Agent id列表
	AgentIdList *[]string `json:"agent_id_list,omitempty"`

	// 开启防护的host id列表
	HostIdList *[]string `json:"host_id_list,omitempty"`
}

func (o ProtectionInfoRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtectionInfoRequestInfo struct{}"
	}

	return strings.Join([]string{"ProtectionInfoRequestInfo", string(data)}, " ")
}
