package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseDocSteps struct {

	// 步骤描述
	ExpectResult *string `json:"expect_result,omitempty"`

	// 预期结果
	TestStep *string `json:"test_step,omitempty"`
}

func (o CaseDocSteps) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseDocSteps struct{}"
	}

	return strings.Join([]string{"CaseDocSteps", string(data)}, " ")
}
