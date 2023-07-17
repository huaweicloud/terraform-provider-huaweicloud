package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateBackupPolicyRequestInfo1 备份策略
type UpdateBackupPolicyRequestInfo1 struct {

	// 策略是否启用，缺省值：true
	Enabled *bool `json:"enabled,omitempty"`

	// 策略ID,若开启防护时开启备份防护，该字段必选
	PolicyId *string `json:"policy_id,omitempty"`

	OperationDefinition *OperationDefinitionRequestInfo `json:"operation_definition,omitempty"`

	Trigger *BackupTriggerRequestInfo1 `json:"trigger,omitempty"`
}

func (o UpdateBackupPolicyRequestInfo1) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBackupPolicyRequestInfo1 struct{}"
	}

	return strings.Join([]string{"UpdateBackupPolicyRequestInfo1", string(data)}, " ")
}
