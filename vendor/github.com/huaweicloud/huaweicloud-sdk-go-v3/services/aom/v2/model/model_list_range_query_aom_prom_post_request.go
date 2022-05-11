package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListRangeQueryAomPromPostRequest struct {

	// PromQL表达式(参考https://prometheus.io/docs/prometheus/latest/querying/basics/)。
	Query string `json:"query"`

	// 起始时间戳(Unix时间戳格式，单位：秒）。
	Start string `json:"start"`

	// 结束时间戳(Unix时间戳格式，单位：秒）。
	End string `json:"end"`

	// 查询时间步长，时间区内每step秒执行一次。
	Step string `json:"step"`
}

func (o ListRangeQueryAomPromPostRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRangeQueryAomPromPostRequest struct{}"
	}

	return strings.Join([]string{"ListRangeQueryAomPromPostRequest", string(data)}, " ")
}
