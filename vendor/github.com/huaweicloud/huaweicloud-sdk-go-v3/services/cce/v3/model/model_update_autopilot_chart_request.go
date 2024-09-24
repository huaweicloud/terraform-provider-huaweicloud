package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAutopilotChartRequest Request Object
type UpdateAutopilotChartRequest struct {

	// 模板的ID
	ChartId string `json:"chart_id"`

	Body *UpdateAutopilotChartRequestBody `json:"body,omitempty" type:"multipart"`
}

func (o UpdateAutopilotChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAutopilotChartRequest struct{}"
	}

	return strings.Join([]string{"UpdateAutopilotChartRequest", string(data)}, " ")
}
