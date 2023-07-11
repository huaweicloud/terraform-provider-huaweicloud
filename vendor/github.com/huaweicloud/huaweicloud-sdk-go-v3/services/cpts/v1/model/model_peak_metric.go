package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PeakMetric struct {

	// vusers
	Vuser *int32 `json:"vuser,omitempty"`

	// tps
	Rps *float64 `json:"rps,omitempty"`

	// avgRT
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
