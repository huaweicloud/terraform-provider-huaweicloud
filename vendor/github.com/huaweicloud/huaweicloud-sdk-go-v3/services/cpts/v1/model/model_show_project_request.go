package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowProjectRequest struct {
	// 测试工程id

	TestSuiteId int32 `json:"test_suite_id"`
}

func (o ShowProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowProjectRequest struct{}"
	}

	return strings.Join([]string{"ShowProjectRequest", string(data)}, " ")
}
