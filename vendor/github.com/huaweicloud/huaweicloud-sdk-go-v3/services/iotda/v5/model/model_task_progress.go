package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 任务执行状态统计对象，用于统计任务的整体执行情况。
type TaskProgress struct {

	// 子任务总个数。
	Total *int32 `json:"total,omitempty"`

	// 正在执行的子任务个数。
	Processing *int32 `json:"processing,omitempty"`

	// 执行成功的子任务个数。
	Success *int32 `json:"success,omitempty"`

	// 执行失败的的子任务个数。
	Fail *int32 `json:"fail,omitempty"`

	// 等待执行的子任务个数。
	Waitting *int32 `json:"waitting,omitempty"`

	// 失败等待重试的子任务个数。
	FailWaitRetry *int32 `json:"fail_wait_retry,omitempty"`

	// 停止的子任务个数。
	Stopped *int32 `json:"stopped,omitempty"`
}

func (o TaskProgress) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskProgress struct{}"
	}

	return strings.Join([]string{"TaskProgress", string(data)}, " ")
}
