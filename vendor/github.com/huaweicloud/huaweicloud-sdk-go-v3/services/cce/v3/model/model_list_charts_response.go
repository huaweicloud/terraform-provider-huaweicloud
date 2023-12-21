package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListChartsResponse Response Object
type ListChartsResponse struct {

	// 模板列表
	Body           *[]ChartResp `json:"body,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListChartsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListChartsResponse struct{}"
	}

	return strings.Join([]string{"ListChartsResponse", string(data)}, " ")
}
