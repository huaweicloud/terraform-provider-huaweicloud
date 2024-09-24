package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotChartsRequest Request Object
type ListAutopilotChartsRequest struct {
}

func (o ListAutopilotChartsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotChartsRequest struct{}"
	}

	return strings.Join([]string{"ListAutopilotChartsRequest", string(data)}, " ")
}
