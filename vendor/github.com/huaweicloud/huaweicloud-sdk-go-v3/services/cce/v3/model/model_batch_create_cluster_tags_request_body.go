package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateClusterTagsRequestBody 批量添加指定集群资源标签的请求体
type BatchCreateClusterTagsRequestBody struct {

	// 待创建的集群资源标签列表。单集群资源标签总数上限为20。
	Tags []ResourceTag `json:"tags"`
}

func (o BatchCreateClusterTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateClusterTagsRequestBody struct{}"
	}

	return strings.Join([]string{"BatchCreateClusterTagsRequestBody", string(data)}, " ")
}
