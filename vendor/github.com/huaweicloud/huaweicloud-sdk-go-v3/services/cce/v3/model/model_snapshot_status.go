package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SnapshotStatus struct {

	// 任务状态
	Phase *string `json:"phase,omitempty"`

	// 任务进度
	Progress *string `json:"progress,omitempty"`

	// 完成时间
	CompletionTime *string `json:"completionTime,omitempty"`
}

func (o SnapshotStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SnapshotStatus struct{}"
	}

	return strings.Join([]string{"SnapshotStatus", string(data)}, " ")
}
