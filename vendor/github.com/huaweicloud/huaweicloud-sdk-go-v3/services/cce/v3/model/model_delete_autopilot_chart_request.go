package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAutopilotChartRequest Request Object
type DeleteAutopilotChartRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`
}

func (o DeleteAutopilotChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAutopilotChartRequest struct{}"
	}

	return strings.Join([]string{"DeleteAutopilotChartRequest", string(data)}, " ")
}
