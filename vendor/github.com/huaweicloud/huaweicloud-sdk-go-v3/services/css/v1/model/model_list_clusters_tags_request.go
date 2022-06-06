package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListClustersTagsRequest struct {

	// 资源类型，当前固定值为“css-cluster”，表示是集群类型。
	ResourceType string `json:"resource_type"`
}

func (o ListClustersTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClustersTagsRequest struct{}"
	}

	return strings.Join([]string{"ListClustersTagsRequest", string(data)}, " ")
}
