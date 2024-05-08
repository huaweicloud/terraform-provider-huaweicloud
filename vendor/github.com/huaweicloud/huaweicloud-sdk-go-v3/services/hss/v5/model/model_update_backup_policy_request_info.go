package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateBackupPolicyRequestInfo 备份策略
type UpdateBackupPolicyRequestInfo struct {

	// 策略是否启用，缺省值：true
	Enabled *bool `json:"enabled,omitempty"`

	// 备份策略ID
	PolicyId string `json:"policy_id"`

	OperationDefinition *OperationDefinitionRequestInfo `json:"operation_definition,omitempty"`

	Trigger *BackupTriggerRequestInfo `json:"trigger,omitempty"`
}

func (o UpdateBackupPolicyRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBackupPolicyRequestInfo struct{}"
	}

	return strings.Join([]string{"UpdateBackupPolicyRequestInfo", string(data)}, " ")
}
