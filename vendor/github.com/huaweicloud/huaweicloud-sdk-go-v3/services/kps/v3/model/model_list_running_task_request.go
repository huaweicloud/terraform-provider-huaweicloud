package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListRunningTaskRequest struct {

	// 每页显示的条目数量。默认值1000。
	Limit *int32 `json:"limit,omitempty"`

	// 首个展示的正在处理任务信息的偏移量
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListRunningTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRunningTaskRequest struct{}"
	}

	return strings.Join([]string{"ListRunningTaskRequest", string(data)}, " ")
}
