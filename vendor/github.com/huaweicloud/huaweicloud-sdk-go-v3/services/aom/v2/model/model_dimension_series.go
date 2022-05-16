package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 时间序列维度列表。
type DimensionSeries struct {

	// 维度名称。
	Name *string `json:"name,omitempty"`

	// 维度取值。
	Value *string `json:"value,omitempty"`
}

func (o DimensionSeries) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DimensionSeries struct{}"
	}

	return strings.Join([]string{"DimensionSeries", string(data)}, " ")
}
