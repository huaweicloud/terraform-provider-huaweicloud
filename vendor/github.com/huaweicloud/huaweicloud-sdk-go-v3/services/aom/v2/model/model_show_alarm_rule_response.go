package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowAlarmRuleResponse struct {
	MetaData *MetaData `json:"meta_data,omitempty"`

	// 阈值规则列表。
	Thresholds     *[]QueryAlarmResult `json:"thresholds,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ShowAlarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAlarmRuleResponse struct{}"
	}

	return strings.Join([]string{"ShowAlarmRuleResponse", string(data)}, " ")
}
