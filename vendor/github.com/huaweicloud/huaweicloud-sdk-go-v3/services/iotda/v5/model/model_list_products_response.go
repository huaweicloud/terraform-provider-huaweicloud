package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListProductsResponse struct {

	// 产品信息列表。
	Products *[]ProductSummary `json:"products,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ListProductsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProductsResponse struct{}"
	}

	return strings.Join([]string{"ListProductsResponse", string(data)}, " ")
}
