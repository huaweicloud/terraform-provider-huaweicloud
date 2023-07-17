package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateReassignmentTaskRequest Request Object
type CreateReassignmentTaskRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *PartitionReassignRequest `json:"body,omitempty"`
}

func (o CreateReassignmentTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateReassignmentTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateReassignmentTaskRequest", string(data)}, " ")
}
