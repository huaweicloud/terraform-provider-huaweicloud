package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMergeTaskCaseRequest Request Object
type ShowMergeTaskCaseRequest struct {

	// 任务运行id（报告id）
	TaskRunId int32 `json:"task_run_id"`
}

func (o ShowMergeTaskCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMergeTaskCaseRequest struct{}"
	}

	return strings.Join([]string{"ShowMergeTaskCaseRequest", string(data)}, " ")
}
