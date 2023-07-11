package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListBackgroundTasksRequest Request Object
type ListBackgroundTasksRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 开启查询的任务编号。
	Start *int32 `json:"start,omitempty"`

	// 查询的任务个数。
	Limit *int32 `json:"limit,omitempty"`

	// 查询任务的最小时间，格式为YYYYMMDDHHmmss。
	BeginTime *string `json:"begin_time,omitempty"`

	// 查询任务的最大时间，格式为YYYYMMDDHHmmss。
	EndTime *string `json:"end_time,omitempty"`
}

func (o ListBackgroundTasksRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListBackgroundTasksRequest struct{}"
	}

	return strings.Join([]string{"ListBackgroundTasksRequest", string(data)}, " ")
}
