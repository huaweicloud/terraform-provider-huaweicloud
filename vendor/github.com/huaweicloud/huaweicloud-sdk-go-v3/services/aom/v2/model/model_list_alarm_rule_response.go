package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListAlarmRuleResponse struct {
	MetaData *MetaData `json:"meta_data,omitempty"`

	Thresholds     *[]QueryAlarmResult `json:"thresholds,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ListAlarmRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAlarmRuleResponse struct{}"
	}

	return strings.Join([]string{"ListAlarmRuleResponse", string(data)}, " ")
}
