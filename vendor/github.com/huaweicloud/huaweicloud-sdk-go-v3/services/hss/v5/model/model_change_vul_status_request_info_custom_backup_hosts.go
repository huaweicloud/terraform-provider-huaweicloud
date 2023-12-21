package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ChangeVulStatusRequestInfoCustomBackupHosts struct {

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 存储库id
	VaultId *string `json:"vault_id,omitempty"`

	// 备份名称
	BackupName *string `json:"backup_name,omitempty"`
}

func (o ChangeVulStatusRequestInfoCustomBackupHosts) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeVulStatusRequestInfoCustomBackupHosts struct{}"
	}

	return strings.Join([]string{"ChangeVulStatusRequestInfoCustomBackupHosts", string(data)}, " ")
}
