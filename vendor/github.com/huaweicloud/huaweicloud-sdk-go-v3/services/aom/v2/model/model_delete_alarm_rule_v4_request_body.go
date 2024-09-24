package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteAlarmRuleV4RequestBody struct {

	// 告警规则名称列表。
	AlarmRules []string `json:"alarm_rules"`
}

func (o DeleteAlarmRuleV4RequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAlarmRuleV4RequestBody struct{}"
	}

	return strings.Join([]string{"DeleteAlarmRuleV4RequestBody", string(data)}, " ")
}
