package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchCreateRequestBody struct {

	// 资源标签
	Tags []BatchCreateRequestBodyTags `json:"tags"`

	// 系统标签
	SysTags *[]BatchCreateRequestBodySysTags `json:"sys_tags,omitempty"`
}

func (o BatchCreateRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateRequestBody struct{}"
	}

	return strings.Join([]string{"BatchCreateRequestBody", string(data)}, " ")
}
