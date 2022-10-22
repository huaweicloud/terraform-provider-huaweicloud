package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteVariableRequest struct {

	// 全局变量id
	VariableId int32 `json:"variable_id"`

	// 工程id
	TestSuiteId int32 `json:"test_suite_id"`
}

func (o DeleteVariableRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteVariableRequest struct{}"
	}

	return strings.Join([]string{"DeleteVariableRequest", string(data)}, " ")
}
