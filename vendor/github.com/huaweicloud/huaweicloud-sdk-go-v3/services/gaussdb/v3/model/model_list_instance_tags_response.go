package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListInstanceTagsResponse struct {
	// 总记录数。

	TotalCount *int32 `json:"total_count,omitempty"`
	// 标签列表。

	Tags           *[]ResourceTagItem `json:"tags,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ListInstanceTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstanceTagsResponse struct{}"
	}

	return strings.Join([]string{"ListInstanceTagsResponse", string(data)}, " ")
}
