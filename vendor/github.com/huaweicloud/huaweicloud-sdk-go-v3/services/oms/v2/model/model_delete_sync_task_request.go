package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSyncTaskRequest Request Object
type DeleteSyncTaskRequest struct {

	// 同步任务ID。
	SyncTaskId string `json:"sync_task_id"`
}

func (o DeleteSyncTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSyncTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteSyncTaskRequest", string(data)}, " ")
}
