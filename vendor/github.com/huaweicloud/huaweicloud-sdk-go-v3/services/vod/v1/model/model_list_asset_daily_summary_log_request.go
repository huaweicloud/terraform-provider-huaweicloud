package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAssetDailySummaryLogRequest Request Object
type ListAssetDailySummaryLogRequest struct {

	// 查询开始时间。仅支持查询一年内的数据，且一次查询的日期跨度不能超过90天。  如果查询指定开始日期的数据，格式为：yyyyMMdd000000。
	StartTime string `json:"start_time"`

	// 查询结束时间。仅支持查询一年内的数据，且一次查询的日期跨度不能超过90天。  如果查询指定结束日期的数据，格式为：yyyyMMdd000000。
	EndTime string `json:"end_time"`

	// 偏移量，表示查询该偏移量后面的记录。
	Offset *int32 `json:"offset,omitempty"`

	// 查询返回记录的数量限制。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListAssetDailySummaryLogRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAssetDailySummaryLogRequest struct{}"
	}

	return strings.Join([]string{"ListAssetDailySummaryLogRequest", string(data)}, " ")
}
