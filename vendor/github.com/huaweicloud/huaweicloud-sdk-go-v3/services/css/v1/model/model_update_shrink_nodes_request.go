package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateShrinkNodesRequest Request Object
type UpdateShrinkNodesRequest struct {

	// 指定待缩容的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *ShrinkNodesReq `json:"body,omitempty"`
}

func (o UpdateShrinkNodesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateShrinkNodesRequest struct{}"
	}

	return strings.Join([]string{"UpdateShrinkNodesRequest", string(data)}, " ")
}
