package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ImageCheckRuleCheckCaseResponseInfo 配置检测检查项的检测用例信息
type ImageCheckRuleCheckCaseResponseInfo struct {

	// 检测用例描述
	CheckDescription *string `json:"check_description,omitempty"`

	// 当前结果
	CurrentValue *string `json:"current_value,omitempty"`

	// 期待结果
	SuggestValue *string `json:"suggest_value,omitempty"`
}

func (o ImageCheckRuleCheckCaseResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImageCheckRuleCheckCaseResponseInfo struct{}"
	}

	return strings.Join([]string{"ImageCheckRuleCheckCaseResponseInfo", string(data)}, " ")
}
