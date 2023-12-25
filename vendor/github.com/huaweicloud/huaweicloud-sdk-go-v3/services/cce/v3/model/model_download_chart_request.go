package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DownloadChartRequest Request Object
type DownloadChartRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`
}

func (o DownloadChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadChartRequest struct{}"
	}

	return strings.Join([]string{"DownloadChartRequest", string(data)}, " ")
}
