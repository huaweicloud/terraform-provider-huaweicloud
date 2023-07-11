package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CoreIndex struct {

	// 错误请求数
	ErrorRequestCount *int32 `json:"error_request_count,omitempty"`

	// 最大RPS
	MaxRps *int32 `json:"max_rps,omitempty"`

	// 最大并发数
	MaxUsers *int32 `json:"max_users,omitempty"`

	// 请求总数
	RequestCount *int32 `json:"request_count,omitempty"`

	// 平均RPS
	Rps *float32 `json:"rps,omitempty"`

	// 成功数
	SuccessCount *int32 `json:"success_count,omitempty"`

	// 成功率
	SuccessRate *int32 `json:"success_rate,omitempty"`

	// 平均TPS
	TransTps *float32 `json:"trans_tps,omitempty"`

	ResponseTime *ResponseTimeInfo `json:"response_time,omitempty"`
}

func (o CoreIndex) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CoreIndex struct{}"
	}

	return strings.Join([]string{"CoreIndex", string(data)}, " ")
}
