package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 修改备份策略信息
type MysqlUpdateBackupPolicyRequest struct {
	BackupPolicy *MysqlBackupPolicy `json:"backup_policy"`
}

func (o MysqlUpdateBackupPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlUpdateBackupPolicyRequest struct{}"
	}

	return strings.Join([]string{"MysqlUpdateBackupPolicyRequest", string(data)}, " ")
}
