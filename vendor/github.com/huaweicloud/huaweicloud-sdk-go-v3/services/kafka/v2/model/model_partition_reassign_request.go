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

	// 是否作为定时任务执行。若非定时执行，则is_schedule 和execute_at字段可为空；若为定时执行，is_schedule为true，execute_at字段非空。
	IsSchedule *bool `json:"is_schedule,omitempty"`

	// 定时时间，格式为Unix时间戳，单位为毫秒
	ExecuteAt *int64 `json:"execute_at,omitempty"`

	// 设为true表示执行时间预估任务，false为执行重平衡任务。
	TimeEstimate *bool `json:"time_estimate,omitempty"`
}

func (o PartitionReassignRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionReassignRequest struct{}"
	}

	return strings.Join([]string{"PartitionReassignRequest", string(data)}, " ")
}
