package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListResourceResponse struct {

	// 资源列表
	Resources *[]Resources `json:"resources,omitempty"`

	// 查询标签下的资源
	Errors *[]Errors `json:"errors,omitempty"`

	// 标签下的资源总数
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListResourceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListResourceResponse struct{}"
	}

	return strings.Join([]string{"ListResourceResponse", string(data)}, " ")
}
