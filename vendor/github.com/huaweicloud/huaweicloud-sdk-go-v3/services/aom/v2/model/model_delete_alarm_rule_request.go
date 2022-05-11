package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteAlarmRuleRequest struct {

	// 阈值规则ID。
	AlarmRuleId string `json:"alarm_rule_id"`
}

func (o DeleteAlarmRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAlarmRuleRequest struct{}"
	}

	return strings.Join([]string{"DeleteAlarmRuleRequest", string(data)}, " ")
}
