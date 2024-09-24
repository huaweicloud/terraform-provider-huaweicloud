package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeletePortTagsRequestBody This is a auto create Body Object
type BatchDeletePortTagsRequestBody struct {

	// 标签列表
	Tags *[]DeleteResourceTagRequestBody `json:"tags,omitempty"`
}

func (o BatchDeletePortTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeletePortTagsRequestBody struct{}"
	}

	return strings.Join([]string{"BatchDeletePortTagsRequestBody", string(data)}, " ")
}
