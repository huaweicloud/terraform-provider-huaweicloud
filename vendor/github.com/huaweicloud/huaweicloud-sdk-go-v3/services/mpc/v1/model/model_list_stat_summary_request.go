package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListStatSummaryRequest Request Object
type ListStatSummaryRequest struct {

	// 开始时间。格式为yyyymmdd。必须是与时区无关的UTC时间。
	StartTime string `json:"start_time"`

	// 结束时间。格式为yyyymmdd。必须是与时区无关的UTC时间。
	EndTime string `json:"end_time"`

	// 支持的参数类型
	StatType string `json:"stat_type"`
}

func (o ListStatSummaryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListStatSummaryRequest struct{}"
	}

	return strings.Join([]string{"ListStatSummaryRequest", string(data)}, " ")
}
