package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateActionRuleResponse Response Object
type UpdateActionRuleResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateActionRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateActionRuleResponse struct{}"
	}

	return strings.Join([]string{"UpdateActionRuleResponse", string(data)}, " ")
}
