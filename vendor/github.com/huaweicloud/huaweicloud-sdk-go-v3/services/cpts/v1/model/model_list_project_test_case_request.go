package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProjectTestCaseRequest Request Object
type ListProjectTestCaseRequest struct {

	// 测试工程id
	TestSuiteId int32 `json:"test_suite_id"`
}

func (o ListProjectTestCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProjectTestCaseRequest struct{}"
	}

	return strings.Join([]string{"ListProjectTestCaseRequest", string(data)}, " ")
}
