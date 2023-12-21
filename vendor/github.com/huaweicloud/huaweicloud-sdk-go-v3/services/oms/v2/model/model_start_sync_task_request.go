package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartSyncTaskRequest Request Object
type StartSyncTaskRequest struct {

	// 同步任务ID。
	SyncTaskId string `json:"sync_task_id"`

	Body *StartSyncTaskReq `json:"body,omitempty"`
}

func (o StartSyncTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartSyncTaskRequest struct{}"
	}

	return strings.Join([]string{"StartSyncTaskRequest", string(data)}, " ")
}
