package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 备份:策略时间调度规则
type BackupTriggerInfo struct {

	// 调度器id
	Id *string `json:"id,omitempty"`

	// 调度器名称
	Name *string `json:"name,omitempty"`

	// 调度器类型,目前只支持 time,定时调度。
	Type *string `json:"type,omitempty"`

	Properties *BackupTriggerPropertiesInfo `json:"properties,omitempty"`
}

func (o BackupTriggerInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupTriggerInfo struct{}"
	}

	return strings.Join([]string{"BackupTriggerInfo", string(data)}, " ")
}
