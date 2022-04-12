package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateProjectRequest struct {
	// 测试工程id

	TestSuiteId int32 `json:"test_suite_id"`

	Body *UpdateProjectRequestBody `json:"body,omitempty"`
}

func (o UpdateProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProjectRequest struct{}"
	}

	return strings.Join([]string{"UpdateProjectRequest", string(data)}, " ")
}
