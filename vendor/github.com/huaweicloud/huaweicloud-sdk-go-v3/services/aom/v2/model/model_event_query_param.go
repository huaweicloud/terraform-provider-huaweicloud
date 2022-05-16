package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 查询事件或者告警信息 。
type EventQueryParam struct {

	// 查询时间范围。 格式：开始时间UTC毫秒.结束时间UTC毫秒.时间范围分钟数。开始和结束时间为-1时，表示最近N分钟，N为时间范围分钟取值。查询时间段，如最近五分钟可以表示为-1.-1.5，固定的时间范围（2017-08-01 08:00:00到2017-08-02 08:00:00）可以表示为1501545600000.1501632000000.1440。
	TimeRange string `json:"time_range"`

	// 统计步长。毫秒数，例如一分钟则填写为60000。
	Step int64 `json:"step"`

	// 模糊查询匹配字段，可以为空。如果值不为空，可以模糊匹配metadata字段中的必选字段的值。
	Search *string `json:"search,omitempty"`

	Sort *EventQueryParamSort `json:"sort,omitempty"`

	// 查询条件组合，可以为空。
	MetadataRelation *[]RelationModel `json:"metadata_relation,omitempty"`
}

func (o EventQueryParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventQueryParam struct{}"
	}

	return strings.Join([]string{"EventQueryParam", string(data)}, " ")
}
