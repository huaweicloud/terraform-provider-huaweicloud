package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProtectionInfoRequestInfo struct {

	// 需要开启防护的主机的操作系统，包含如下：   - Windows : Windows系统   - Linux : Linux系统
	OperatingSystem string `json:"operating_system"`

	// 勒索防护是否开启，包含如下：   - closed ：关闭。   - opened ：开启。   若选择开启，protection_policy_id或者create_protection_policy必填一项
	RansomProtectionStatus string `json:"ransom_protection_status"`

	// 勒索防护策略ID,若选择已有策略防护,则该字段必选
	ProtectionPolicyId *string `json:"protection_policy_id,omitempty"`

	CreateProtectionPolicy *ProtectionProxyInfoRequestInfo `json:"create_protection_policy,omitempty"`

	// 是否服务器备份，包含如下：   - closed ：关闭。   - opened ：开启。   若选择开启服务器备份，则backup_cycle必填
	BackupProtectionStatus string `json:"backup_protection_status"`

	BackupResources *BackupResources `json:"backup_resources,omitempty"`

	// 备份策略ID
	BackupPolicyId *string `json:"backup_policy_id,omitempty"`

	BackupCycle *UpdateBackupPolicyRequestInfo1 `json:"backup_cycle,omitempty"`

	// 开启防护的Agent id列表
	AgentIdList []string `json:"agent_id_list"`

	// 开启防护的host id列表
	HostIdList []string `json:"host_id_list"`
}

func (o ProtectionInfoRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtectionInfoRequestInfo struct{}"
	}

	return strings.Join([]string{"ProtectionInfoRequestInfo", string(data)}, " ")
}
