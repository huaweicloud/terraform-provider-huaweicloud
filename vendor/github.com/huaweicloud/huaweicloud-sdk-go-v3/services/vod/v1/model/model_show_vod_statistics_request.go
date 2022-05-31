package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowVodStatisticsRequest struct {

	// 起始时间，格式为yyyymmddhhmmss。
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间，格式为yyyymmddhhmmss。 - “start_time”、“end_time”均不存在时，“start_time”取当天零点，“end_time”取当前时间。 - “start_time”不存在、“end_time”存在，请求非法。 - “start_time”存在、“end_time”不存在，“end_time”取当前时间。 - 只能查询最近三个月内的数据，且时间跨度不能超过31天。 - 起始时间和结束时间会自动规整，起始时间规整为指定时间所在的整点时刻，结束时间规整为指定时间所在时间的下一小时整点时刻。
	EndTime *string `json:"end_time,omitempty"`

	// 查询粒度间隔。  取值如下： - 时间跨度1天：1小时、4小时、8小时，分别对应3600秒、14400秒和28800秒。 - 时间跨度2~7天：1小时、4小时、8小时、1天，分别对应3600秒、14400秒、28800秒和86400秒。 - 时间跨度8~31天：4小时、8小时、1天，分别对应14400秒、28800秒和86400秒。  单位：秒。  若不设置，默认取对应时间跨度的最小间隔。
	Interval *int32 `json:"interval,omitempty"`
}

func (o ShowVodStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVodStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ShowVodStatisticsRequest", string(data)}, " ")
}
