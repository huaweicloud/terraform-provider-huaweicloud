package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteEvent2alarmRuleResponse Response Object
type DeleteEvent2alarmRuleResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteEvent2alarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteEvent2alarmRuleResponse struct{}"
	}

	return strings.Join([]string{"DeleteEvent2alarmRuleResponse", string(data)}, " ")
}
