package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchDeleteResourceTagsRequestBody struct {

	// 标签列表。
	Tags *[]Tags `json:"tags,omitempty"`
}

func (o BatchDeleteResourceTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteResourceTagsRequestBody struct{}"
	}

	return strings.Join([]string{"BatchDeleteResourceTagsRequestBody", string(data)}, " ")
}
