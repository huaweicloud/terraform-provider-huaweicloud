package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListLiveSampleLogsRequest struct {

	// 播放域名。
	PlayDomain string `json:"play_domain"`

	// 查询开始时间，UTC时间：YYYY-MM-DDTHH:mm:ssZ，如北京时间2020年3月4日16点00分00秒可表示为2020-03-04T08:00:00Z。仅支持查询最近3个月内的数据。
	StartTime *sdktime.SdkTime `json:"start_time"`

	// 查询结束时间，UTC时间：YYYY-MM-DDTHH:mm:ssZ，如北京时间2020年3月4日16点00分00秒可表示为2020-03-04T08:00:00Z。查询时间跨度不能大于7天。
	EndTime *sdktime.SdkTime `json:"end_time"`
}

func (o ListLiveSampleLogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLiveSampleLogsRequest struct{}"
	}

	return strings.Join([]string{"ListLiveSampleLogsRequest", string(data)}, " ")
}
