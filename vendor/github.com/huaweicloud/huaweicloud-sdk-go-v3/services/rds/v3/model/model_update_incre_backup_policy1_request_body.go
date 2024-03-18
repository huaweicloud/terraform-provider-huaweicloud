package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateIncreBackupPolicy1RequestBody update incre backup policy
type UpdateIncreBackupPolicy1RequestBody struct {
	IncreBackupPolicy *ShowIncreBackupPolicyRespBodyIncreBackupPolicy `json:"incre_backup_policy"`
}

func (o UpdateIncreBackupPolicy1RequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateIncreBackupPolicy1RequestBody struct{}"
	}

	return strings.Join([]string{"UpdateIncreBackupPolicy1RequestBody", string(data)}, " ")
}
