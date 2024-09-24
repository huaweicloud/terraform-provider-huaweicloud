package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddOrUpdateMetricOrEventAlarmRuleResponse Response Object
type AddOrUpdateMetricOrEventAlarmRuleResponse struct {

	// 错误码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 错误信息。
	ErrorMessage *string `json:"error_message,omitempty"`

	// 告警规则列表。
	AlarmRules     *[]AddOrUpdateAlarmRuleV4ItemResult `json:"alarm_rules,omitempty"`
	HttpStatusCode int                                 `json:"-"`
}

func (o AddOrUpdateMetricOrEventAlarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddOrUpdateMetricOrEventAlarmRuleResponse struct{}"
	}

	return strings.Join([]string{"AddOrUpdateMetricOrEventAlarmRuleResponse", string(data)}, " ")
}
