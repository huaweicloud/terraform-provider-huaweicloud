package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 按标签查询资源请求结构体。
type QueryResourceByTagsDto struct {

	// **参数说明**：要查询的资源类型，当前支持设备（device）。
	ResourceType string `json:"resource_type"`

	// **参数说明**：标签列表，支持按照标签key和value组合查询，传入的多个标签之间是或的关系。
	Tags []TagV5Dto `json:"tags"`
}

func (o QueryResourceByTagsDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueryResourceByTagsDto struct{}"
	}

	return strings.Join([]string{"QueryResourceByTagsDto", string(data)}, " ")
}
