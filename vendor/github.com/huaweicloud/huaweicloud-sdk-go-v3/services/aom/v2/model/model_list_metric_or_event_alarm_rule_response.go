package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListMetricOrEventAlarmRuleResponse Response Object
type ListMetricOrEventAlarmRuleResponse struct {

	// 告警规则列表。
	AlarmRules *[]AlarmParamForV4Db `json:"alarm_rules,omitempty"`

	// 元数据信息。
	Metadata *interface{} `json:"metadata,omitempty"`

	// 告警规则数量。
	Count          *int32 `json:"count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListMetricOrEventAlarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMetricOrEventAlarmRuleResponse struct{}"
	}

	return strings.Join([]string{"ListMetricOrEventAlarmRuleResponse", string(data)}, " ")
}
