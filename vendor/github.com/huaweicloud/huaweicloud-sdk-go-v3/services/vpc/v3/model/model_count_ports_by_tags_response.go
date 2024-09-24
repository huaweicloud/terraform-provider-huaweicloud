package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CountPortsByTagsResponse Response Object
type CountPortsByTagsResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	// 资源数量
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CountPortsByTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CountPortsByTagsResponse struct{}"
	}

	return strings.Join([]string{"CountPortsByTagsResponse", string(data)}, " ")
}
