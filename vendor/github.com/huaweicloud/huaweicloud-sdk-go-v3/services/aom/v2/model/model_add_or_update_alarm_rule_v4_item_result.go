package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddOrUpdateAlarmRuleV4ItemResult struct {

	// 告警规则名称。
	AlarmRuleName string `json:"alarm_rule_name"`

	// 告警规则新增或修改结果。
	Result string `json:"result"`
}

func (o AddOrUpdateAlarmRuleV4ItemResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddOrUpdateAlarmRuleV4ItemResult struct{}"
	}

	return strings.Join([]string{"AddOrUpdateAlarmRuleV4ItemResult", string(data)}, " ")
}
