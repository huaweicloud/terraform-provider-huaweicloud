package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMergeReportLogsOutlineRequest Request Object
type ShowMergeReportLogsOutlineRequest struct {

	// 任务运行id（报告id）
	TaskRunId int32 `json:"task_run_id"`
}

func (o ShowMergeReportLogsOutlineRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMergeReportLogsOutlineRequest struct{}"
	}

	return strings.Join([]string{"ShowMergeReportLogsOutlineRequest", string(data)}, " ")
}
