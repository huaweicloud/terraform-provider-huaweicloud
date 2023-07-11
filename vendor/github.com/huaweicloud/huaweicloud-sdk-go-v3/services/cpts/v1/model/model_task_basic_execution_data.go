package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TaskBasicExecutionData struct {

	// 已执行完成的用例数
	CompleteNum *int32 `json:"complete_num,omitempty"`

	// 持续时间
	Duration *int32 `json:"duration,omitempty"`

	// 结束时间
	EndTime *string `json:"end_time,omitempty"`

	// 已执行用例数
	ExecutedNum *int32 `json:"executed_num,omitempty"`

	// 【指标数据:最后一个轮次】 用例数
	KpiCaseCount *int32 `json:"kpi_case_count,omitempty"`

	// 【指标数据:最后一个轮次】 已执行的用例数
	KpiCaseExecuteCount *int32 `json:"kpi_case_execute_count,omitempty"`

	// 【指标数据:最后一个轮次】 最后一轮结果为Pass的用例数
	KpiCasePassCount *int32 `json:"kpi_case_pass_count,omitempty"`

	// 任务间用例是否并行执行
	Parallel *bool `json:"parallel,omitempty"`

	// 用例通过数
	PassNum *int32 `json:"pass_num,omitempty"`

	// 开始时间
	StartTime *string `json:"start_time,omitempty"`

	// 任务状态
	TaskStatus *int32 `json:"task_status,omitempty"`

	// 总用例数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 分钟*并发数
	Vum *int32 `json:"vum,omitempty"`
}

func (o TaskBasicExecutionData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskBasicExecutionData struct{}"
	}

	return strings.Join([]string{"TaskBasicExecutionData", string(data)}, " ")
}
