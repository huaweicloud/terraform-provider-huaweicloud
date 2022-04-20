package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DebugCaseRequest struct {
	// 测试工程id

	TestSuiteId int32 `json:"test_suite_id"`
	// 任务id

	TaskId int32 `json:"task_id"`
	// 用例id

	CaseId int32 `json:"case_id"`

	Body *DebugCaseRequestBody `json:"body,omitempty"`
}

func (o DebugCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseRequest struct{}"
	}

	return strings.Join([]string{"DebugCaseRequest", string(data)}, " ")
}
