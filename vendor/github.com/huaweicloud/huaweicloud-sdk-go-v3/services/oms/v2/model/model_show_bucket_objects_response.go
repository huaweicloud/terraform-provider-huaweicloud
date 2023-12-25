package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBucketObjectsResponse Response Object
type ShowBucketObjectsResponse struct {

	// 是否存在下一页
	NextPage *bool `json:"next_page,omitempty"`

	// 查询下一页所需要的标记（此页末尾对象名，偏移量）
	NextMarker *string `json:"next_marker,omitempty"`

	// 查询桶信息返回的record
	Records        *[]ShowBucketRecord `json:"records,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ShowBucketObjectsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBucketObjectsResponse struct{}"
	}

	return strings.Join([]string{"ShowBucketObjectsResponse", string(data)}, " ")
}
