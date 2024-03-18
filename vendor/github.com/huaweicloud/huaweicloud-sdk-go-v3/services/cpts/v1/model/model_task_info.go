package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TaskInfo struct {

	// 基准并发
	BenchConcurrent *int32 `json:"bench_concurrent,omitempty"`

	// 用例Id列表
	CaseList *[]CaseInfoDetail `json:"case_list,omitempty"`

	// 创建时间
	CreateTime *sdktime.SdkTime `json:"create_time,omitempty"`

	// 描述信息
	Description *string `json:"description,omitempty"`

	// 任务名称
	Name *string `json:"name,omitempty"`

	// 任务模式（0：时长模式；1：次数模式；2：混合模式）
	OperateMode *int32 `json:"operate_mode,omitempty"`

	// 任务所属工程id
	ProjectId *int32 `json:"project_id,omitempty"`

	// 最近一次运行的报告简略信息，包括运行任务id，即本对象的task_run_info_id。运行用例id，即本对象的related_temp_running_id。
	RelatedTempRunningData *[]RelatedTempRunningData `json:"related_temp_running_data,omitempty"`

	// 任务运行状态（9：等待运行；0：运行中；1：暂停；2：结束； 3：异常中止；4：用户主动终止（完成状态）；5：用户主动终止）
	RunStatus *int32 `json:"run_status,omitempty"`

	// 任务更新时间
	UpdateTime *string `json:"update_time,omitempty"`

	// 任务间用例是否并行执行
	Parallel *bool `json:"parallel,omitempty"`
}

func (o TaskInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskInfo struct{}"
	}

	return strings.Join([]string{"TaskInfo", string(data)}, " ")
}
