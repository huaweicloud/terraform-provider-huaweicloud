package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateBackupPolicyInfoResponse Response Object
type UpdateBackupPolicyInfoResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateBackupPolicyInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBackupPolicyInfoResponse struct{}"
	}

	return strings.Join([]string{"UpdateBackupPolicyInfoResponse", string(data)}, " ")
}
