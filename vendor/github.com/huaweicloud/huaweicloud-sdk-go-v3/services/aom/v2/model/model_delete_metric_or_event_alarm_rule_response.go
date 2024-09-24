package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteMetricOrEventAlarmRuleResponse Response Object
type DeleteMetricOrEventAlarmRuleResponse struct {

	// 错误码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 错误信息。
	ErrorMessage *string `json:"error_message,omitempty"`

	// 资源列表。
	Resources      *[]ItemResult `json:"resources,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o DeleteMetricOrEventAlarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMetricOrEventAlarmRuleResponse struct{}"
	}

	return strings.Join([]string{"DeleteMetricOrEventAlarmRuleResponse", string(data)}, " ")
}
