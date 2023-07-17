package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateEventRuleRequest Request Object
type UpdateEventRuleRequest struct {
	Body *Event2alarmRuleBody `json:"body,omitempty"`
}

func (o UpdateEventRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateEventRuleRequest struct{}"
	}

	return strings.Join([]string{"UpdateEventRuleRequest", string(data)}, " ")
}
