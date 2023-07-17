package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResponseTimeInfo struct {

	// 平均响应时间
	AvgResponseTime *float32 `json:"avg_response_time,omitempty"`

	// TP50
	AvgTp50 *int32 `json:"avg_tp50,omitempty"`

	// TP90
	AvgTp90 *int32 `json:"avg_tp90,omitempty"`

	// 最大响应时间
	MaxResponseTime *int32 `json:"max_response_time,omitempty"`

	// 最小响应时间
	MinResponseTime *int32 `json:"min_response_time,omitempty"`
}

func (o ResponseTimeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResponseTimeInfo struct{}"
	}

	return strings.Join([]string{"ResponseTimeInfo", string(data)}, " ")
}
