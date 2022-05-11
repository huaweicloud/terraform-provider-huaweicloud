package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 时序数据查询参数。
type QuerySampleParam struct {

	// 时序数据对象列表。  取值范围：JSON数组大小不超过20。
	Samples []QuerySample `json:"samples"`

	// 统计方式。 取值范围： maximum，minimum，sum，average，sampleCount。
	Statistics []string `json:"statistics"`

	// 监控数据粒度。 取值范围 枚举值，取值范围： 60，1分钟粒度； 300，5分钟粒度； 900，15分钟粒度； 3600，1小时粒度。
	Period int32 `json:"period"`

	// 说明： time_range/period≤1440 计算时，time_range和period需换算为相同的单位。 取值范围 格式：开始时间UTC毫秒.结束时间UTC毫秒.时间范围分钟数。开始和结束时间为-1时，表示最近N分钟，N为时间范围分钟取值。 查询时间段，如最近五分钟可以表示为-1.-1.5，固定的时间范围（2017-08-01 08:00 :00到2017-08-02 08:00:00）可以表示为1501545600000.1501632000000.1440。
	TimeRange string `json:"time_range"`
}

func (o QuerySampleParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QuerySampleParam struct{}"
	}

	return strings.Join([]string{"QuerySampleParam", string(data)}, " ")
}
