package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteMetricOrEventAlarmRuleRequest Request Object
type DeleteMetricOrEventAlarmRuleRequest struct {
	Body *DeleteAlarmRuleV4RequestBody `json:"body,omitempty"`
}

func (o DeleteMetricOrEventAlarmRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMetricOrEventAlarmRuleRequest struct{}"
	}

	return strings.Join([]string{"DeleteMetricOrEventAlarmRuleRequest", string(data)}, " ")
}
