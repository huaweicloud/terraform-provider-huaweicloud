package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TaskInfo struct {
	// bench_concurrent

	BenchConcurrent *int32 `json:"bench_concurrent,omitempty"`
	// case_list

	CaseList *[]CaseInfo `json:"case_list,omitempty"`
	// 创建时间

	CreateTime *sdktime.SdkTime `json:"create_time,omitempty"`
	// description

	Description *string `json:"description,omitempty"`
	// name

	Name *string `json:"name,omitempty"`
	// operate_mode

	OperateMode *int32 `json:"operate_mode,omitempty"`
	// project_id

	ProjectId *int32 `json:"project_id,omitempty"`
	// related_temp_running_data

	RelatedTempRunningData *[]RelatedTempRunningData `json:"related_temp_running_data,omitempty"`
	// run_status

	RunStatus *int32 `json:"run_status,omitempty"`
	// update_time

	UpdateTime *string `json:"update_time,omitempty"`
	// parallel

	Parallel *bool `json:"parallel,omitempty"`
}

func (o TaskInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskInfo struct{}"
	}

	return strings.Join([]string{"TaskInfo", string(data)}, " ")
}
