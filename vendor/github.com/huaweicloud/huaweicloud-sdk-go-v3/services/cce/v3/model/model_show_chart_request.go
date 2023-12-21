package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowChartRequest Request Object
type ShowChartRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`
}

func (o ShowChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowChartRequest struct{}"
	}

	return strings.Join([]string{"ShowChartRequest", string(data)}, " ")
}
