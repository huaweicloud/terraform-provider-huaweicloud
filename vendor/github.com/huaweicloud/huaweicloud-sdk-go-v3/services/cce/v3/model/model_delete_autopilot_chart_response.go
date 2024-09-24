package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteAutopilotChartResponse Response Object
type DeleteAutopilotChartResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteAutopilotChartResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteAutopilotChartResponse struct{}"
	}

	return strings.Join([]string{"DeleteAutopilotChartResponse", string(data)}, " ")
}
