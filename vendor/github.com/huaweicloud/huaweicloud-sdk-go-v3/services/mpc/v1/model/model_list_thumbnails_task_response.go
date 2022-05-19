package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListThumbnailsTaskResponse struct {

	// 截图任务组
	TaskArray *[]ThumbTask `json:"task_array,omitempty"`

	// 查询结果是否被截取。 - 1代表被截取，即还有结果未被返回，可以通过设置page和size参数继续查询。 - 0代表未被截取，即所有结果已被返回。
	IsTruncated *int32 `json:"is_truncated,omitempty"`

	// 查询结果的数量。
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListThumbnailsTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListThumbnailsTaskResponse struct{}"
	}

	return strings.Join([]string{"ListThumbnailsTaskResponse", string(data)}, " ")
}
