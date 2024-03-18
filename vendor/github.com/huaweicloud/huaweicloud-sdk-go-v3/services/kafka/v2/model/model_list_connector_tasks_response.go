package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListConnectorTasksResponse Response Object
type ListConnectorTasksResponse struct {

	// Smart Connector任务详情。
	Tasks *[]SmartConnectTaskEntity `json:"tasks,omitempty"`

	// Smart Connector任务数。
	TotalNumber *int32 `json:"total_number,omitempty"`

	// Smart Connector最大任务数。
	MaxTasks *int32 `json:"max_tasks,omitempty"`

	// Smart Connector任务配额。
	QuotaTasks     *int32 `json:"quota_tasks,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListConnectorTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListConnectorTasksResponse struct{}"
	}

	return strings.Join([]string{"ListConnectorTasksResponse", string(data)}, " ")
}
