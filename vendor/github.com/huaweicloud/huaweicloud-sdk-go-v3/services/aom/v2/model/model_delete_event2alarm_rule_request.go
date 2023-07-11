package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteEvent2alarmRuleRequest Request Object
type DeleteEvent2alarmRuleRequest struct {
	Body *[]string `json:"body,omitempty"`
}

func (o DeleteEvent2alarmRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteEvent2alarmRuleRequest struct{}"
	}

	return strings.Join([]string{"DeleteEvent2alarmRuleRequest", string(data)}, " ")
}
