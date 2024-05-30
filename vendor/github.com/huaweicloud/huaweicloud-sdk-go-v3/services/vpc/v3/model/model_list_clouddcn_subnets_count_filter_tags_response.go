package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClouddcnSubnetsCountFilterTagsResponse Response Object
type ListClouddcnSubnetsCountFilterTagsResponse struct {

	// 本次请求的编号
	RequestId *string `json:"request_id,omitempty"`

	// 当前列表中资源数量。
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListClouddcnSubnetsCountFilterTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClouddcnSubnetsCountFilterTagsResponse struct{}"
	}

	return strings.Join([]string{"ListClouddcnSubnetsCountFilterTagsResponse", string(data)}, " ")
}
