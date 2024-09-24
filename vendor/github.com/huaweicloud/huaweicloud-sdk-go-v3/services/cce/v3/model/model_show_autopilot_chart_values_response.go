package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotChartValuesResponse Response Object
type ShowAutopilotChartValuesResponse struct {
	Values         *ChartValueValues `json:"values,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ShowAutopilotChartValuesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotChartValuesResponse struct{}"
	}

	return strings.Join([]string{"ShowAutopilotChartValuesResponse", string(data)}, " ")
}
