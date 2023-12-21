package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteChartResponse Response Object
type DeleteChartResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteChartResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteChartResponse struct{}"
	}

	return strings.Join([]string{"DeleteChartResponse", string(data)}, " ")
}
