package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteActionRuleResponse Response Object
type DeleteActionRuleResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteActionRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteActionRuleResponse struct{}"
	}

	return strings.Join([]string{"DeleteActionRuleResponse", string(data)}, " ")
}
