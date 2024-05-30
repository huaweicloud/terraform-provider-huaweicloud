package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClouddcnSubnetsFilterTagsResponse Response Object
type ListClouddcnSubnetsFilterTagsResponse struct {

	// 资源列表
	Resources *[]ClouddcnResource `json:"resources,omitempty"`

	// 当前列表中资源数量。
	TotalCount *int32 `json:"total_count,omitempty"`

	// 本次请求的编号
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListClouddcnSubnetsFilterTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClouddcnSubnetsFilterTagsResponse struct{}"
	}

	return strings.Join([]string{"ListClouddcnSubnetsFilterTagsResponse", string(data)}, " ")
}
