package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteAlarmRuleResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteAlarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAlarmRuleResponse struct{}"
	}

	return strings.Join([]string{"DeleteAlarmRuleResponse", string(data)}, " ")
}
