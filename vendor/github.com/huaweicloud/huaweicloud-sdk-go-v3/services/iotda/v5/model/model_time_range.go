package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则有效时间段
type TimeRange struct {

	// **参数说明**：规则条件触发的开始时间，格式：HH:mm。
	StartTime string `json:"start_time"`

	// **参数说明**：规则条件触发的结束时间，格式：HH:mm。若结束时间与开始时间一致，则时间为全天。
	EndTime string `json:"end_time"`

	// **参数说明**：星期列表，以逗号分隔。1代表周日，2代表周一，依次类推，默认为每天。星期列表中的日期为开始时间的日期。 **取值范围**：只允许数字和逗号的组合，数字不小于1不大于7，数量不超过7个，以逗号隔开
	DaysOfWeek *string `json:"days_of_week,omitempty"`
}

func (o TimeRange) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TimeRange struct{}"
	}

	return strings.Join([]string{"TimeRange", string(data)}, " ")
}
