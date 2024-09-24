package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotChartsResponse Response Object
type ListAutopilotChartsResponse struct {

	// 模板列表
	Body           *[]ChartResp `json:"body,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListAutopilotChartsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotChartsResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotChartsResponse", string(data)}, " ")
}
