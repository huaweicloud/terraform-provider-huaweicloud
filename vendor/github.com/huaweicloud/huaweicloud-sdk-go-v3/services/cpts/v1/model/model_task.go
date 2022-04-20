package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// tasks
type Task struct {
	// bench_concurrent

	BenchConcurrent *int32 `json:"bench_concurrent,omitempty"`
	// description

	Description *string `json:"description,omitempty"`
	// id

	Id *int32 `json:"id,omitempty"`
	// name

	Name *string `json:"name,omitempty"`
	// operate_mode

	OperateMode *int32 `json:"operate_mode,omitempty"`

	TaskRunInfo *TaskRunInfo `json:"task_run_info,omitempty"`
	// update_time

	UpdateTime *string `json:"update_time,omitempty"`
}

func (o Task) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Task struct{}"
	}

	return strings.Join([]string{"Task", string(data)}, " ")
}
