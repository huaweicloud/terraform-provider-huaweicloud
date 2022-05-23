package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 单个子任务详情结构体
type TaskDetail struct {

	// 执行批量任务的目标。
	Target *string `json:"target,omitempty"`

	// 子任务的执行状态，结果范围：Processing，Success，Fail，Waitting，FailWaitRetry，Stopped。 - Waitting: 等待执行。 - Processing: 执行中。 - Success: 成功。 - Fail: 失败。 - FailWaitRetry: 失败重试。 - Stopped: 已停止。
	Status *string `json:"status,omitempty"`

	// 子任务执行的输出信息。
	Output *string `json:"output,omitempty"`

	Error *ErrorInfo `json:"error,omitempty"`
}

func (o TaskDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskDetail struct{}"
	}

	return strings.Join([]string{"TaskDetail", string(data)}, " ")
}
