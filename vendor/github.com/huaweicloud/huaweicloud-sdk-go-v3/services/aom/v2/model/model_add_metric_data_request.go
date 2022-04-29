package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type AddMetricDataRequest struct {
	Body *[]MetricDataItem `json:"body,omitempty"`
}

func (o AddMetricDataRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddMetricDataRequest struct{}"
	}

	return strings.Join([]string{"AddMetricDataRequest", string(data)}, " ")
}
