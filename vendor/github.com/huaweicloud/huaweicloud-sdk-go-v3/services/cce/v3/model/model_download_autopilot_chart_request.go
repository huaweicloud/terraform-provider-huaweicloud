package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DownloadAutopilotChartRequest Request Object
type DownloadAutopilotChartRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`
}

func (o DownloadAutopilotChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadAutopilotChartRequest struct{}"
	}

	return strings.Join([]string{"DownloadAutopilotChartRequest", string(data)}, " ")
}
