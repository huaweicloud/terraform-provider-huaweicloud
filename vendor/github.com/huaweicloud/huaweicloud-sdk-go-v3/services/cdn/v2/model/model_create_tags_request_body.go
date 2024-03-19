package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTagsRequestBody CreateTagsRequestBody
type CreateTagsRequestBody struct {

	// 资源id。  > 域名ID
	ResourceId string `json:"resource_id"`

	// 标签列表。
	Tags []TagMap `json:"tags"`
}

func (o CreateTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTagsRequestBody struct{}"
	}

	return strings.Join([]string{"CreateTagsRequestBody", string(data)}, " ")
}
