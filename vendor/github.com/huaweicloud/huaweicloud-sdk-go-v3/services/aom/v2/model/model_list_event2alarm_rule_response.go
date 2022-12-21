package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListEvent2alarmRuleResponse struct {
	Body           *[]Event2alarmRuleBody `json:"body,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ListEvent2alarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEvent2alarmRuleResponse struct{}"
	}

	return strings.Join([]string{"ListEvent2alarmRuleResponse", string(data)}, " ")
}
