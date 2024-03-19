package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeCheckRuleActionResponse Response Object
type ChangeCheckRuleActionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeCheckRuleActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeCheckRuleActionResponse struct{}"
	}

	return strings.Join([]string{"ChangeCheckRuleActionResponse", string(data)}, " ")
}
