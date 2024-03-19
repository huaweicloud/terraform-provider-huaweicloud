package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTaskRequestBody UpdateTaskRequestBody
type UpdateTaskRequestBody struct {

	// 任务id
	Id int32 `json:"id"`

	// 任务名称
	Name string `json:"name"`

	// 描述
	Description *string `json:"description,omitempty"`

	// 工程id
	ProjectId int32 `json:"project_id"`

	// 任务运行状态（9：等待运行；0：运行中；1：暂停；2：结束； 3：异常中止；4：用户主动终止（完成状态）；5：用户主动终止）
	RunStatus *int32 `json:"run_status,omitempty"`

	// 任务类型（0：旧版本任务；1：融合版本任务）
	RunType *int32 `json:"run_type,omitempty"`

	TaskRunInfo *TaskRunInfo `json:"task_run_info,omitempty"`

	// 用例信息
	CaseList *[]CaseInfoDetail `json:"case_list,omitempty"`

	// 压力阶段模式，0：时长模式；1：次数模式
	OperateMode *int32 `json:"operate_mode,omitempty"`

	// 基准并发
	BenchConcurrent *int32 `json:"bench_concurrent,omitempty"`

	// 最近一次运行的报告简略信息
	RelatedTempRunningData *[]RelatedTempRunningData `json:"related_temp_running_data,omitempty"`
}

func (o UpdateTaskRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateTaskRequestBody", string(data)}, " ")
}
