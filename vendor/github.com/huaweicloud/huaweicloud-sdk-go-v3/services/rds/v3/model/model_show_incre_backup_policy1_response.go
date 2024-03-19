package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowIncreBackupPolicy1Response Response Object
type ShowIncreBackupPolicy1Response struct {
	IncreBackupPolicy *ShowIncreBackupPolicyRespBodyIncreBackupPolicy `json:"incre_backup_policy,omitempty"`
	HttpStatusCode    int                                             `json:"-"`
}

func (o ShowIncreBackupPolicy1Response) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowIncreBackupPolicy1Response struct{}"
	}

	return strings.Join([]string{"ShowIncreBackupPolicy1Response", string(data)}, " ")
}
