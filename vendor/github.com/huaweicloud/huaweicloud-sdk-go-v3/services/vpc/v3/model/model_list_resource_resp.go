package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListResourceResp
type ListResourceResp struct {

	// 资源ID
	ResourceId string `json:"resource_id"`

	// 资源详情。 资源对象，用于扩展。默认为空
	ResourceDetail *interface{} `json:"resource_detail"`

	// 资源名称，资源没有名称时默认为空字符串.
	ResourceName string `json:"resource_name"`

	// 标签列表，没有标签默认为空数组
	Tags []ResourceTag `json:"tags"`
}

func (o ListResourceResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListResourceResp struct{}"
	}

	return strings.Join([]string{"ListResourceResp", string(data)}, " ")
}
