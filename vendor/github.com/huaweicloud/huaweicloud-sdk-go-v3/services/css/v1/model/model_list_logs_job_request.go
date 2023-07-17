package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListLogsJobRequest Request Object
type ListLogsJobRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`

	// 指定查询起始值，默认值为1，即从第1个任务开始查询。
	Start *int32 `json:"start,omitempty"`

	// 指定查询个数，默认值为10，即一次查询10个任务信息。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListLogsJobRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLogsJobRequest struct{}"
	}

	return strings.Join([]string{"ListLogsJobRequest", string(data)}, " ")
}
