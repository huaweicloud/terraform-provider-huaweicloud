package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 带宽限速策略。
type BandwidthPolicyDto struct {

	// 流量控制开始时间（包含），格式为“hh:mm”。例如“12:03”表示12时03分。
	End string `json:"end"`

	// 时段内允许的最大流量带宽，单位Byte/s，取值范围为>= 1048576Byte/s（相当于1MB/s）且<=209715200Byte/s（相当于200MB/s）。
	MaxBandwidth int64 `json:"max_bandwidth"`

	// 流量控制开始时间（包含），格式为“hh:mm”。例如“12:03”表示12时03分。
	Start string `json:"start"`
}

func (o BandwidthPolicyDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BandwidthPolicyDto struct{}"
	}

	return strings.Join([]string{"BandwidthPolicyDto", string(data)}, " ")
}
