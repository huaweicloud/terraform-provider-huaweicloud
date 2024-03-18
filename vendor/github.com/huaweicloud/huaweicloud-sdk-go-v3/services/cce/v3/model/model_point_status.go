package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PointStatus struct {
	TaskType *TaskType `json:"taskType,omitempty"`

	// 升级任务项ID
	TaskID *string `json:"taskID,omitempty"`

	Status *UpgradeWorkflowTaskStatus `json:"status,omitempty"`

	// 升级任务开始时间
	StartTimeStamp *string `json:"startTimeStamp,omitempty"`

	// 升级任务结束时间
	EndTimeStamp *string `json:"endTimeStamp,omitempty"`

	// 升级任务过期时间（当前仅升级前检查任务适用）
	ExpireTimeStamp *string `json:"expireTimeStamp,omitempty"`
}

func (o PointStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PointStatus struct{}"
	}

	return strings.Join([]string{"PointStatus", string(data)}, " ")
}
