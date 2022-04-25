package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListAlarmRuleRequest struct {

	// 分页信息。
	Offset *string `json:"offset,omitempty"`

	// 用于限制结果数据条数。 取值范围(0,1000],默认值为1000。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListAlarmRuleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAlarmRuleRequest struct{}"
	}

	return strings.Join([]string{"ListAlarmRuleRequest", string(data)}, " ")
}
