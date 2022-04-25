package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowMetricsDataResponse struct {

	// 响应码。
	ErrorCode *string `json:"errorCode,omitempty"`

	// 响应信息描述。
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// 指标对象列表。
	Metrics        *[]MetricDataValue `json:"metrics,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ShowMetricsDataResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMetricsDataResponse struct{}"
	}

	return strings.Join([]string{"ShowMetricsDataResponse", string(data)}, " ")
}
