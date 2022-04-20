package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateTaskStatusRequest struct {
	// 测试工程id

	TestSuiteId int32 `json:"test_suite_id"`
	// 任务id

	TaskId int32 `json:"task_id"`

	Body *UpdateTaskStatusRequestBody `json:"body,omitempty"`
}

func (o UpdateTaskStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskStatusRequest struct{}"
	}

	return strings.Join([]string{"UpdateTaskStatusRequest", string(data)}, " ")
}
