package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PrecheckTaskMetadata 升级前检查任务元数据
type PrecheckTaskMetadata struct {

	// 任务ID
	Uid *string `json:"uid,omitempty"`

	// 任务创建时间
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`

	// 任务更新时间
	UpdateTimestamp *string `json:"updateTimestamp,omitempty"`
}

func (o PrecheckTaskMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrecheckTaskMetadata struct{}"
	}

	return strings.Join([]string{"PrecheckTaskMetadata", string(data)}, " ")
}
