package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSyncTasksResponse Response Object
type ListSyncTasksResponse struct {

	// 查询的同步任务详情
	Tasks *[]SyncTaskInfo `json:"tasks,omitempty"`

	// 满足查询条件的同步任务总数
	Count          *int64 `json:"count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListSyncTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSyncTasksResponse struct{}"
	}

	return strings.Join([]string{"ListSyncTasksResponse", string(data)}, " ")
}
