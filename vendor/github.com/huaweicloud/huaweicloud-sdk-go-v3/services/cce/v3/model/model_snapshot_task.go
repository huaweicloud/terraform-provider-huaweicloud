package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SnapshotTask struct {

	// 任务类型
	Kind *string `json:"kind,omitempty"`

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *SnapshotTaskMetadata `json:"metadata,omitempty"`

	Spec *SnapshotSpec `json:"spec,omitempty"`

	Status *SnapshotStatus `json:"status,omitempty"`
}

func (o SnapshotTask) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SnapshotTask struct{}"
	}

	return strings.Join([]string{"SnapshotTask", string(data)}, " ")
}
