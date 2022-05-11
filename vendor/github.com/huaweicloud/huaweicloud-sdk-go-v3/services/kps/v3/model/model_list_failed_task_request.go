package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListFailedTaskRequest struct {

	// 每页返回的个数。 默认值：50。
	Limit *string `json:"limit,omitempty"`

	// 偏移量，表示从此偏移量开始查询， offset大于等于0
	Offset *string `json:"offset,omitempty"`
}

func (o ListFailedTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFailedTaskRequest struct{}"
	}

	return strings.Join([]string{"ListFailedTaskRequest", string(data)}, " ")
}
