package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckRuleCheckCaseResponseInfo 配置检测检查项的检测用例信息
type CheckRuleCheckCaseResponseInfo struct {

	// 检测用例描述
	CheckDescription *string `json:"check_description,omitempty"`

	// 当前结果
	CurrentValue *string `json:"current_value,omitempty"`

	// 期待结果
	SuggestValue *string `json:"suggest_value,omitempty"`
}

func (o CheckRuleCheckCaseResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckRuleCheckCaseResponseInfo struct{}"
	}

	return strings.Join([]string{"CheckRuleCheckCaseResponseInfo", string(data)}, " ")
}
