package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClustersTagsRequest Request Object
type CreateClustersTagsRequest struct {

	// 指定待添加标签的集群ID。
	ClusterId string `json:"cluster_id"`

	// 资源类型，当前固定值为“css-cluster”，表示是集群类型。
	ResourceType string `json:"resource_type"`

	Body *TagReq `json:"body,omitempty"`
}

func (o CreateClustersTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClustersTagsRequest struct{}"
	}

	return strings.Join([]string{"CreateClustersTagsRequest", string(data)}, " ")
}
