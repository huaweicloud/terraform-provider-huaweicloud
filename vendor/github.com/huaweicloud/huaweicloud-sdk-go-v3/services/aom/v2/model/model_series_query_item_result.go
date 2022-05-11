package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 维度信息。
type SeriesQueryItemResult struct {

	// 命名空间。
	Namespace *string `json:"namespace,omitempty"`

	// 维度列表。
	Dimensions *[]DimensionSeries `json:"dimensions,omitempty"`

	// 时间序列名称。
	MetricName *string `json:"metric_name,omitempty"`

	// 时间序列单位。
	Unit *string `json:"unit,omitempty"`

	// 时间序列哈希值。
	DimensionValueHash *string `json:"dimension_value_hash,omitempty"`
}

func (o SeriesQueryItemResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SeriesQueryItemResult struct{}"
	}

	return strings.Join([]string{"SeriesQueryItemResult", string(data)}, " ")
}
