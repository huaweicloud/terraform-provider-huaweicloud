package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreaseCaseResponseJson 响应json
type CreaseCaseResponseJson struct {

	// 用例id
	TestCaseId *int32 `json:"test_case_id,omitempty"`
}

func (o CreaseCaseResponseJson) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreaseCaseResponseJson struct{}"
	}

	return strings.Join([]string{"CreaseCaseResponseJson", string(data)}, " ")
}
