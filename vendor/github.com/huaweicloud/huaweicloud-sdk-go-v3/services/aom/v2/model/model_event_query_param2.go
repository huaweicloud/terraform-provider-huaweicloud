package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EventQueryParam2 查询事件或者告警信息 。
type EventQueryParam2 struct {

	// timeRange用于指标查询时间范围，主要用于解决客户端时间和服务端时间不一致情况下，查询最近N分钟的数据。另可用于精确查询某一段时间的数据。   如：   - -1.-1.60(表示最近60分钟)，不管当前客户端是什么时间，都以服务端时间为准查询最近60分钟。  - 1650852000000.1650852300000.5(表示GMT+8 2022-04-25 10:00:00至2022-04-25 10:05:00指定的5分钟)   格式：   startTimeInMillis.endTimeInMillis.durationInMinutes   参数解释：   - startTimeInMillis: 查询的开始时间，格式为UTC毫秒，如果指定为-1，服务端将按(endTimeInMillis - durationInMinutes * 60 * 1000)计算开始时间。如-1.1650852300000.5，则相当于1650852000000.1650852300000.5  - endTimeInMillis: 查询的结束时间，格式为UTC毫秒，如果指定为-1，服务端将按(startTimeInMillis + durationInMinutes * 60 * 1000)计算结束时间，如果计算出的结束时间大于当前系统时间，则使用当前系统时间。如1650852000000.-1.5，则相当于1650852000000.1650852300000.5  - durationInMinutes：查询时间的跨度分钟数。 取值范围大于0并且大于等于(endTimeInMillis - startTimeInMillis) / (60 * 1000) - 1。当开始时间与结束时间都设置为-1时，系统会将结束时间设置为当前时间UTC毫秒值，并按(endTimeInMillis - durationInMinutes * 60 * 1000)计算开始时间。如：-1.-1.60(表示最近60分钟)   约束：   单次请求中，查询时长与周期需要满足以下条件: durationInMinutes * 60 / period <= 1440
	TimeRange string `json:"time_range"`

	// 统计步长。毫秒数，例如一分钟则填写为60000。
	Step *int64 `json:"step,omitempty"`

	// 模糊查询匹配字段，可以为空。如果值不为空，可以模糊匹配。metadata字段为必选字段。
	Search *string `json:"search,omitempty"`

	Sort *EventQueryParam2Sort `json:"sort,omitempty"`

	// 查询条件组合，可以为空。
	MetadataRelation *[]RelationModel `json:"metadata_relation,omitempty"`
}

func (o EventQueryParam2) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventQueryParam2 struct{}"
	}

	return strings.Join([]string{"EventQueryParam2", string(data)}, " ")
}
