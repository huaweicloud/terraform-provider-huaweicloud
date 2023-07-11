package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDirectoryRequest Request Object
type UpdateDirectoryRequest struct {

	// 目录id
	DirectoryId int32 `json:"directory_id"`

	// 测试工程id
	TestSuiteId int32 `json:"test_suite_id"`

	Body *UpdateDirectoryRequestBody `json:"body,omitempty"`
}

func (o UpdateDirectoryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDirectoryRequest struct{}"
	}

	return strings.Join([]string{"UpdateDirectoryRequest", string(data)}, " ")
}
