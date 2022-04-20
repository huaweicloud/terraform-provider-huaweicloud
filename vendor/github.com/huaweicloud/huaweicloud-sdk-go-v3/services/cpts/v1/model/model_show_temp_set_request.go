package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowTempSetRequest struct {
	// 测试工程id

	TestSuiteId int32 `json:"test_suite_id"`
	// 查询偏移

	Offset *int32 `json:"offset,omitempty"`
	// 查询数量

	Limit *int32 `json:"limit,omitempty"`
}

func (o ShowTempSetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTempSetRequest struct{}"
	}

	return strings.Join([]string{"ShowTempSetRequest", string(data)}, " ")
}
