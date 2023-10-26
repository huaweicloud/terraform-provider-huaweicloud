package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopHotPipelineRequest Request Object
type StopHotPipelineRequest struct {

	// 指定待操作的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *StopHotPipelineRequestBody `json:"body,omitempty"`
}

func (o StopHotPipelineRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopHotPipelineRequest struct{}"
	}

	return strings.Join([]string{"StopHotPipelineRequest", string(data)}, " ")
}
