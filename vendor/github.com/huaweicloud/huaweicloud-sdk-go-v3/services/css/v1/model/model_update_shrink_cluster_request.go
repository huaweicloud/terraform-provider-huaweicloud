package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateShrinkClusterRequest struct {

	// 指定待缩容的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *ShrinkClusterReq `json:"body,omitempty"`
}

func (o UpdateShrinkClusterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateShrinkClusterRequest struct{}"
	}

	return strings.Join([]string{"UpdateShrinkClusterRequest", string(data)}, " ")
}
