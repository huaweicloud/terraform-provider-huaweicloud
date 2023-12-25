package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowChartValuesResponse Response Object
type ShowChartValuesResponse struct {
	Values         *ChartValueValues `json:"values,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ShowChartValuesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowChartValuesResponse struct{}"
	}

	return strings.Join([]string{"ShowChartValuesResponse", string(data)}, " ")
}
