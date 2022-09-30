package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteClustersTagsRequest struct {

	// 指定删除标签的集群ID。
	ClusterId string `json:"cluster_id"`

	// 资源类型，当前固定值为“css-cluster”，表示是集群类型。
	ResourceType string `json:"resource_type"`

	// 需要删除的标签名。
	Key string `json:"key"`
}

func (o DeleteClustersTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteClustersTagsRequest struct{}"
	}

	return strings.Join([]string{"DeleteClustersTagsRequest", string(data)}, " ")
}
