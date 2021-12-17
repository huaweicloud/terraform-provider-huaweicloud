package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTasksResponse struct {
	// 符合要求的任务数量，不受分页影响

	Count *int32 `json:"count,omitempty"`
	// 查询到的任务列表

	Tasks          *[]TasksResponseBody `json:"tasks,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ListTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTasksResponse struct{}"
	}

	return strings.Join([]string{"ListTasksResponse", string(data)}, " ")
}
