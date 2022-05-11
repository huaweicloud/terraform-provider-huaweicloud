package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteAlarmRulesBody struct {

	// 阈值规则列表
	AlarmRules []string `json:"alarm_rules"`
}

func (o DeleteAlarmRulesBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAlarmRulesBody struct{}"
	}

	return strings.Join([]string{"DeleteAlarmRulesBody", string(data)}, " ")
}
