package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRunningTaskResponse struct {

	// 正在处理的任务总数。
	Total *int32 `json:"total,omitempty"`

	// 正在处理的任务列表。
	Tasks          *[]RunningTasks `json:"tasks,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ListRunningTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRunningTaskResponse struct{}"
	}

	return strings.Join([]string{"ListRunningTaskResponse", string(data)}, " ")
}
