package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotChartRequest Request Object
type ShowAutopilotChartRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`
}

func (o ShowAutopilotChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotChartRequest struct{}"
	}

	return strings.Join([]string{"ShowAutopilotChartRequest", string(data)}, " ")
}
