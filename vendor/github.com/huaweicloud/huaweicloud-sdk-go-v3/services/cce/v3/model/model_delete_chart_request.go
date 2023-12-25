package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteChartRequest Request Object
type DeleteChartRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`
}

func (o DeleteChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteChartRequest struct{}"
	}

	return strings.Join([]string{"DeleteChartRequest", string(data)}, " ")
}
