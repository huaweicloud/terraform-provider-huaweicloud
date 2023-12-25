package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowSyncTaskRequest Request Object
type ShowSyncTaskRequest struct {

	// 同步任务ID。
	SyncTaskId string `json:"sync_task_id"`

	// 查询同步任务详情的时间（毫秒），依据该值返回所在月份的统计数据。
	QueryTime string `json:"query_time"`
}

func (o ShowSyncTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowSyncTaskRequest struct{}"
	}

	return strings.Join([]string{"ShowSyncTaskRequest", string(data)}, " ")
}
