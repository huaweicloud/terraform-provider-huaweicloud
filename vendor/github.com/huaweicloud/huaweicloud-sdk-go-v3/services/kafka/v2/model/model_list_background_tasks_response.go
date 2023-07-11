package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListBackgroundTasksResponse Response Object
type ListBackgroundTasksResponse struct {

	// 任务数量。
	TaskCount *string `json:"task_count,omitempty"`

	// 任务列表。
	Tasks          *[]ListBackgroundTasksRespTasks `json:"tasks,omitempty"`
	HttpStatusCode int                             `json:"-"`
}

func (o ListBackgroundTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListBackgroundTasksResponse struct{}"
	}

	return strings.Join([]string{"ListBackgroundTasksResponse", string(data)}, " ")
}
