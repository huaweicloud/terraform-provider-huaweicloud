package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PartitionReassignRequest struct {

	// 重平衡分配方案。
	Reassignments []PartitionReassignEntity `json:"reassignments"`

	// 重平衡门限值。
	Throttle *int32 `json:"throttle,omitempty"`
}

func (o PartitionReassignRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionReassignRequest struct{}"
	}

	return strings.Join([]string{"PartitionReassignRequest", string(data)}, " ")
}
