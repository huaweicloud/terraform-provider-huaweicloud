package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ItemResult struct {

	// 告警规则名称。
	AlarmRuleName string `json:"alarm_rule_name"`
}

func (o ItemResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ItemResult struct{}"
	}

	return strings.Join([]string{"ItemResult", string(data)}, " ")
}
