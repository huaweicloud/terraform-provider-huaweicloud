package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchDeleteRequestBody struct {

	// 资源标签
	Tags []BatchDeleteRequestBodyTags `json:"tags"`

	// 系统标签
	SysTags *[]BatchDeleteRequestBodySysTags `json:"sys_tags,omitempty"`
}

func (o BatchDeleteRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteRequestBody struct{}"
	}

	return strings.Join([]string{"BatchDeleteRequestBody", string(data)}, " ")
}
