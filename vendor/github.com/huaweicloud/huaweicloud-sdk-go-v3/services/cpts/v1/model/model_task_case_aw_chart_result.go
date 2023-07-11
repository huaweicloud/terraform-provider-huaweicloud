package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TaskCaseAwChartResult struct {
	BrokenList *BrokenList `json:"broken_list,omitempty"`

	// 错误信息
	ErrMessage *string `json:"err_message,omitempty"`

	// 响应时间区间与出现次数的汇总信息
	RespTimeRange map[string]string `json:"resp_time_range,omitempty"`
}

func (o TaskCaseAwChartResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskCaseAwChartResult struct{}"
	}

	return strings.Join([]string{"TaskCaseAwChartResult", string(data)}, " ")
}
