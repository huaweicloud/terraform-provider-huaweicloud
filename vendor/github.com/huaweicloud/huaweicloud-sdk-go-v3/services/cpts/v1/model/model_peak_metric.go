package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PeakMetric struct {

	// 虚拟用户数
	Vuser *int32 `json:"vuser,omitempty"`

	// 每秒事务数
	Rps *float64 `json:"rps,omitempty"`

	// 平均响应时间
	AvgRT *float64 `json:"avgRT,omitempty"`

	// 成功率
	SuccessRate *float64 `json:"successRate,omitempty"`

	// 峰值时间
	PeakTime *string `json:"peakTime,omitempty"`
}

func (o PeakMetric) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PeakMetric struct{}"
	}

	return strings.Join([]string{"PeakMetric", string(data)}, " ")
}
