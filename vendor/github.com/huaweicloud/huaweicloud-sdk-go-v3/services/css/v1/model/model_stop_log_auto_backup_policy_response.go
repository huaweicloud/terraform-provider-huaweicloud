package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopLogAutoBackupPolicyResponse Response Object
type StopLogAutoBackupPolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopLogAutoBackupPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopLogAutoBackupPolicyResponse struct{}"
	}

	return strings.Join([]string{"StopLogAutoBackupPolicyResponse", string(data)}, " ")
}
