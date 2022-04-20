package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowReportRequest struct {
	// 运行任务id

	TaskRunId int32 `json:"task_run_id"`
	// 运行用例id

	CaseRunId int32 `json:"case_run_id"`
	// 曲线图点数

	BrokensLimitCount int32 `json:"brokens_limit_count"`
}

func (o ShowReportRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowReportRequest struct{}"
	}

	return strings.Join([]string{"ShowReportRequest", string(data)}, " ")
}
