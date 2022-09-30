package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTagsRequestBody
type DeleteTagsRequestBody struct {

	// 资源id
	ResourceId string `json:"resource_id"`

	// 键列表
	Tags []string `json:"tags"`
}

func (o DeleteTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTagsRequestBody struct{}"
	}

	return strings.Join([]string{"DeleteTagsRequestBody", string(data)}, " ")
}
