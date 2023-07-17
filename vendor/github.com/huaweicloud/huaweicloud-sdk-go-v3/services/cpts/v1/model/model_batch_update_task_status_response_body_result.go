package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchUpdateTaskStatusResponseBodyResult struct {

	// 任务id
	TaskId *int32 `json:"task_id,omitempty"`

	// 报告id
	TaskRunId *int32 `json:"task_run_id,omitempty"`
}

func (o BatchUpdateTaskStatusResponseBodyResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchUpdateTaskStatusResponseBodyResult struct{}"
	}

	return strings.Join([]string{"BatchUpdateTaskStatusResponseBodyResult", string(data)}, " ")
}
