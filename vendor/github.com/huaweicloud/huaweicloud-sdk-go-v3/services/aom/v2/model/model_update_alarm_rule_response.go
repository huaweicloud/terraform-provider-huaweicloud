package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateAlarmRuleResponse struct {

	// 阈值规则id。
	AlarmRuleId    *int64 `json:"alarm_rule_id,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o UpdateAlarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAlarmRuleResponse struct{}"
	}

	return strings.Join([]string{"UpdateAlarmRuleResponse", string(data)}, " ")
}
