package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateBatchClustersTagsRequest Request Object
type UpdateBatchClustersTagsRequest struct {

	// 指定批量添加标签的集群ID。
	ClusterId string `json:"cluster_id"`

	// 资源类型，当前固定值为“css-cluster”，表示是集群类型。
	ResourceType string `json:"resource_type"`

	Body *BatchAddOrDeleteTagOnClusterReq `json:"body,omitempty"`
}

func (o UpdateBatchClustersTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBatchClustersTagsRequest struct{}"
	}

	return strings.Join([]string{"UpdateBatchClustersTagsRequest", string(data)}, " ")
}
