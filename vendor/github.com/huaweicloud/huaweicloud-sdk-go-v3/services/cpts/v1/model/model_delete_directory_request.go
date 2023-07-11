package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDirectoryRequest Request Object
type DeleteDirectoryRequest struct {

	// 目录id
	DirectoryId int32 `json:"directory_id"`

	// 测试工程id
	TestSuiteId int32 `json:"test_suite_id"`
}

func (o DeleteDirectoryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDirectoryRequest struct{}"
	}

	return strings.Join([]string{"DeleteDirectoryRequest", string(data)}, " ")
}
