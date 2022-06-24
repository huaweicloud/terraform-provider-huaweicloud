package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateSyncEventsRequest struct {

	// 同步任务ID
	SyncTaskId string `json:"sync_task_id"`

	Body *SyncObjectReq `json:"body,omitempty"`
}

func (o CreateSyncEventsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSyncEventsRequest struct{}"
	}

	return strings.Join([]string{"CreateSyncEventsRequest", string(data)}, " ")
}
