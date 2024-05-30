package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClouddcnSubnetsTagsResponse Response Object
type ListClouddcnSubnetsTagsResponse struct {

	// 本次请求的编号
	RequestId *string `json:"request_id,omitempty"`

	// 当前列表中资源数量。
	TotalCount *int32 `json:"total_count,omitempty"`

	// tag列表信息
	Tags           *[]ResourceTags `json:"tags,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ListClouddcnSubnetsTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClouddcnSubnetsTagsResponse struct{}"
	}

	return strings.Join([]string{"ListClouddcnSubnetsTagsResponse", string(data)}, " ")
}
