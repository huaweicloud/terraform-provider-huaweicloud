package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPortsByTagsResponse Response Object
type ListPortsByTagsResponse struct {

	// 资源列表
	Resources *[]ListResourceResp `json:"resources,omitempty"`

	// 资源数量
	TotalCount *int32 `json:"total_count,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListPortsByTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPortsByTagsResponse struct{}"
	}

	return strings.Join([]string{"ListPortsByTagsResponse", string(data)}, " ")
}
