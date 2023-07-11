package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListEvent2alarmRuleRequest Request Object
type ListEvent2alarmRuleRequest struct {
}

func (o ListEvent2alarmRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEvent2alarmRuleRequest struct{}"
	}

	return strings.Join([]string{"ListEvent2alarmRuleRequest", string(data)}, " ")
}
