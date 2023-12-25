package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSyncTaskResponse Response Object
type CreateSyncTaskResponse struct {

	// 同步任务ID
	SyncTaskId     *string `json:"sync_task_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateSyncTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSyncTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateSyncTaskResponse", string(data)}, " ")
}
