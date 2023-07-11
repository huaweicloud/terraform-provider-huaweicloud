package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BackupTriggerRequestInfo1 策略时间调度规则
type BackupTriggerRequestInfo1 struct {
	Properties *BackupTriggerPropertiesRequestInfo1 `json:"properties,omitempty"`
}

func (o BackupTriggerRequestInfo1) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupTriggerRequestInfo1 struct{}"
	}

	return strings.Join([]string{"BackupTriggerRequestInfo1", string(data)}, " ")
}
