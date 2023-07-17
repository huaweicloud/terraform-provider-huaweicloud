package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddEvent2alarmRuleRequest Request Object
type AddEvent2alarmRuleRequest struct {
	Body *Event2alarmRuleBody `json:"body,omitempty"`
}

func (o AddEvent2alarmRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddEvent2alarmRuleRequest struct{}"
	}

	return strings.Join([]string{"AddEvent2alarmRuleRequest", string(data)}, " ")
}
