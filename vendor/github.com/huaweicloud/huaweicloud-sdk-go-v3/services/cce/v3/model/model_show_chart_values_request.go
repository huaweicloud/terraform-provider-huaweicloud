package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowChartValuesRequest Request Object
type ShowChartValuesRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`
}

func (o ShowChartValuesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowChartValuesRequest struct{}"
	}

	return strings.Join([]string{"ShowChartValuesRequest", string(data)}, " ")
}
