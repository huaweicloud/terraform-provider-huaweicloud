package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type AddMetricDataResponse struct {

	// 响应码。
	ErrorCode *string `json:"errorCode,omitempty"`

	// 响应信息描述。
	ErrorMessage   *string `json:"errorMessage,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AddMetricDataResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddMetricDataResponse struct{}"
	}

	return strings.Join([]string{"AddMetricDataResponse", string(data)}, " ")
}
