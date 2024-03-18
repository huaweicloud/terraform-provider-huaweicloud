package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowHistoryTasksResponse Response Object
type ShowHistoryTasksResponse struct {

	// 总共的任务个数。
	Total *int32 `json:"total,omitempty"`

	// 日志列表数据
	Tasks *[]TasksObject `json:"tasks,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowHistoryTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowHistoryTasksResponse struct{}"
	}

	return strings.Join([]string{"ShowHistoryTasksResponse", string(data)}, " ")
}
