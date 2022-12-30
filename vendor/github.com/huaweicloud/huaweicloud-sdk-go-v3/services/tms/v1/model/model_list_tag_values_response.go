package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTagValuesResponse struct {

	// 查询到的标签值列表
	Values *[]string `json:"values,omitempty"`

	PageInfo       *PageInfoTagValues `json:"page_info,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ListTagValuesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTagValuesResponse struct{}"
	}

	return strings.Join([]string{"ListTagValuesResponse", string(data)}, " ")
}
