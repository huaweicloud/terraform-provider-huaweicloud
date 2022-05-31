package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListAssetCategoryResponse struct {

	// 分类返回值
	Body           *[]QueryCategoryRsp `json:"body,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ListAssetCategoryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAssetCategoryResponse struct{}"
	}

	return strings.Join([]string{"ListAssetCategoryResponse", string(data)}, " ")
}
