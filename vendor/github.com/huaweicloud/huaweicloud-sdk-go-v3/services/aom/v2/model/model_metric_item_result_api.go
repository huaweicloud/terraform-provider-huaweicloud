package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 维度信息。
type MetricItemResultApi struct {

	// 指标维度列表。
	Dimensions *[]Dimension `json:"dimensions,omitempty"`

	// 指标哈希值。
	Dimensionvaluehash *string `json:"dimensionvaluehash,omitempty"`

	// 指标名称。
	MetricName *string `json:"metricName,omitempty"`

	// 命名空间。
	Namespace *string `json:"namespace,omitempty"`

	// 指标单位。
	Unit *string `json:"unit,omitempty"`
}

func (o MetricItemResultApi) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MetricItemResultApi struct{}"
	}

	return strings.Join([]string{"MetricItemResultApi", string(data)}, " ")
}
