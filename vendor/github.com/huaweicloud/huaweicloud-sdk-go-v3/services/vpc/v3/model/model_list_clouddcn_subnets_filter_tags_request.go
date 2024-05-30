package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClouddcnSubnetsFilterTagsRequest Request Object
type ListClouddcnSubnetsFilterTagsRequest struct {

	// 每页返回的个数
	Limit *int32 `json:"limit,omitempty"`

	// 分页起始点
	Offset *int32 `json:"offset,omitempty"`

	Body *ListResourcesByTagsRequestBody `json:"body,omitempty"`
}

func (o ListClouddcnSubnetsFilterTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClouddcnSubnetsFilterTagsRequest struct{}"
	}

	return strings.Join([]string{"ListClouddcnSubnetsFilterTagsRequest", string(data)}, " ")
}
