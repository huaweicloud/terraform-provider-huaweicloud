package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CountEventsResponse Response Object
type CountEventsResponse struct {

	// 统计步长。毫秒数，例如一分钟则填写为60000。
	Step *int64 `json:"step,omitempty"`

	// 统计结果对应的时间序列。
	Timestamps *[]int64 `json:"timestamps,omitempty"`

	// 事件或者告警不同级别相同时间序列对应的统计结果。
	Series *[]EventSeries `json:"series,omitempty"`

	// 各类告警信息的数量汇总
	Summary        map[string]int64 `json:"summary,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o CountEventsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CountEventsResponse struct{}"
	}

	return strings.Join([]string{"CountEventsResponse", string(data)}, " ")
}
