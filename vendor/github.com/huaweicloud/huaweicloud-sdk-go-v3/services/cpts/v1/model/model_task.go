package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// tasks
type Task struct {

	// 基准并发
	BenchConcurrent *int32 `json:"bench_concurrent,omitempty"`

	// 描述信息
	Description *string `json:"description,omitempty"`

	// 任务Id
	Id *int32 `json:"id,omitempty"`

	// 任务名称
	Name *string `json:"name,omitempty"`

	// 任务压测模式
	OperateMode *int32 `json:"operate_mode,omitempty"`

	TaskRunInfo *TaskRunInfo `json:"task_run_info,omitempty"`

	// 更新时间
	UpdateTime *string `json:"update_time,omitempty"`

	// 任务间用例是否并行执行
	Parallel *bool `json:"parallel,omitempty"`
}

func (o Task) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Task struct{}"
	}

	return strings.Join([]string{"Task", string(data)}, " ")
}
