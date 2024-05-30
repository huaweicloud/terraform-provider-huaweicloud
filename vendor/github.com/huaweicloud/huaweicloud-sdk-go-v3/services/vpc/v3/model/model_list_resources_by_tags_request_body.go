package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ListResourcesByTagsRequestBody struct {

	// 包含标签。
	Tags *[]ResourceTags `json:"tags,omitempty"`
}

func (o ListResourcesByTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListResourcesByTagsRequestBody struct{}"
	}

	return strings.Join([]string{"ListResourcesByTagsRequestBody", string(data)}, " ")
}
