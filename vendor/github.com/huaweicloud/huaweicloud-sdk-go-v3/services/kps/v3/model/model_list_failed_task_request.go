package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListFailedTaskRequest Request Object
type ListFailedTaskRequest struct {

	// 每页显示的条目数量。默认值1000。
	Limit *int32 `json:"limit,omitempty"`

	// 失败的任务信息列表的偏移量
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListFailedTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFailedTaskRequest struct{}"
	}

	return strings.Join([]string{"ListFailedTaskRequest", string(data)}, " ")
}
