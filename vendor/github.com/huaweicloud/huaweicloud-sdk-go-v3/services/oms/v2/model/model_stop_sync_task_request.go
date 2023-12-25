package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopSyncTaskRequest Request Object
type StopSyncTaskRequest struct {

	// 同步任务ID。
	SyncTaskId string `json:"sync_task_id"`
}

func (o StopSyncTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopSyncTaskRequest struct{}"
	}

	return strings.Join([]string{"StopSyncTaskRequest", string(data)}, " ")
}
