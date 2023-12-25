package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChartValueValues values.yaml中的数据，数据结构以具体的模板为准
type ChartValueValues struct {
	Basic *interface{} `json:"basic,omitempty"`
}

func (o ChartValueValues) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChartValueValues struct{}"
	}

	return strings.Join([]string{"ChartValueValues", string(data)}, " ")
}
