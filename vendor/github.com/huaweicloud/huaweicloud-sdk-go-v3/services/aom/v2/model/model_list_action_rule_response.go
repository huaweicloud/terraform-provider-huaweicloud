package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListActionRuleResponse Response Object
type ListActionRuleResponse struct {

	// 告警行动规则列表
	ActionRules    *[]ActionRule `json:"action_rules,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o ListActionRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListActionRuleResponse struct{}"
	}

	return strings.Join([]string{"ListActionRuleResponse", string(data)}, " ")
}
