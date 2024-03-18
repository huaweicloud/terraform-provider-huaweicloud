package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckRuleFixValuesInfo 用户键入的要修复的检查项的参数ID和参数值
type CheckRuleFixValuesInfo struct {

	// 检查项的参数ID
	RuleParamId *int32 `json:"rule_param_id,omitempty"`

	// 检查项的参数值
	FixValue *int32 `json:"fix_value,omitempty"`
}

func (o CheckRuleFixValuesInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckRuleFixValuesInfo struct{}"
	}

	return strings.Join([]string{"CheckRuleFixValuesInfo", string(data)}, " ")
}
