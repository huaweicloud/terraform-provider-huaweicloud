package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTranscodingTaskResponse struct {

	// 返回任务组
	TaskArray *[]QueryTranscodingsTaskResponse `json:"task_array,omitempty"`

	// 查询结果是否被截取。 - 1：表示被截取，即还有结果未被返回，可以通过设置page和size参数继续查询。 - 0：表示未被截取，即所有结果已被返回。
	IsTruncated *int32 `json:"is_truncated,omitempty"`

	// 查询结果的数量。
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListTranscodingTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTranscodingTaskResponse struct{}"
	}

	return strings.Join([]string{"ListTranscodingTaskResponse", string(data)}, " ")
}
