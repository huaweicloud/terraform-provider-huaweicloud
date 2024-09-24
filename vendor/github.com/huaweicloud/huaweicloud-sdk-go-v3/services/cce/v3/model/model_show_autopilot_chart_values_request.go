package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotChartValuesRequest Request Object
type ShowAutopilotChartValuesRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`
}

func (o ShowAutopilotChartValuesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotChartValuesRequest struct{}"
	}

	return strings.Join([]string{"ShowAutopilotChartValuesRequest", string(data)}, " ")
}
