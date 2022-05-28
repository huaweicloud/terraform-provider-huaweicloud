package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowCdnStatisticsResponse struct {

	// 统计起始时间。
	StartTime *string `json:"start_time,omitempty"`

	// 统计间隔。
	Interval *int32 `json:"interval,omitempty"`

	// 采样数据数组。  从start_time开始，每个间隔对应一个采样数据。
	Values         *[]int64 `json:"values,omitempty"`
	HttpStatusCode int      `json:"-"`
}

func (o ShowCdnStatisticsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCdnStatisticsResponse struct{}"
	}

	return strings.Join([]string{"ShowCdnStatisticsResponse", string(data)}, " ")
}
