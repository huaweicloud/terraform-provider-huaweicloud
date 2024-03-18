package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMergeReportLogsOutlineRequest Request Object
type ShowMergeReportLogsOutlineRequest struct {

	// 运行任务id，即报告id。启动任务（更新任务状态或批量启停任务）接口，会返回运行任务id。
	TaskRunId int32 `json:"task_run_id"`
}

func (o ShowMergeReportLogsOutlineRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMergeReportLogsOutlineRequest struct{}"
	}

	return strings.Join([]string{"ShowMergeReportLogsOutlineRequest", string(data)}, " ")
}
