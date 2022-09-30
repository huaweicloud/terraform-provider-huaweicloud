package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type StopPipelineRequest struct {

	// 指定停止pipeline的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o StopPipelineRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopPipelineRequest struct{}"
	}

	return strings.Join([]string{"StopPipelineRequest", string(data)}, " ")
}
