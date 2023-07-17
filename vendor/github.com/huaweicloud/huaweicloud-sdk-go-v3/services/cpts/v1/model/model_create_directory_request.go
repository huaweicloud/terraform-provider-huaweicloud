package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateDirectoryRequest Request Object
type CreateDirectoryRequest struct {

	// 测试工程id
	TestSuiteId int32 `json:"test_suite_id"`

	Body *CreateDirectoryRequestBody `json:"body,omitempty"`
}

func (o CreateDirectoryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDirectoryRequest struct{}"
	}

	return strings.Join([]string{"CreateDirectoryRequest", string(data)}, " ")
}
