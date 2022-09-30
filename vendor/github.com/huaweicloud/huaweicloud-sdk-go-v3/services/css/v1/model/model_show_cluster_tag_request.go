package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowClusterTagRequest struct {

	// 指定待查询的集群ID。
	ClusterId string `json:"cluster_id"`

	// 资源类型，当前固定值为“css-cluster”，表示是集群类型。
	ResourceType string `json:"resource_type"`
}

func (o ShowClusterTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterTagRequest struct{}"
	}

	return strings.Join([]string{"ShowClusterTagRequest", string(data)}, " ")
}
