package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTaskCaseAwChartRequest Request Object
type ShowTaskCaseAwChartRequest struct {

	// 任务运行id（报告id）
	TaskRunId int32 `json:"task_run_id"`

	// 用例运行id
	CaseRunId int32 `json:"case_run_id"`

	// 详情id
	DetailId string `json:"detail_id"`
}

func (o ShowTaskCaseAwChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskCaseAwChartRequest struct{}"
	}

	return strings.Join([]string{"ShowTaskCaseAwChartRequest", string(data)}, " ")
}
