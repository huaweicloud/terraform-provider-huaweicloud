package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListProjectTagsResponse struct {
	// 总记录数。

	TotalCount *int32 `json:"total_count,omitempty"`
	// 标签列表。

	Tags           *[]ProjectTagItem `json:"tags,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ListProjectTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProjectTagsResponse struct{}"
	}

	return strings.Join([]string{"ListProjectTagsResponse", string(data)}, " ")
}
