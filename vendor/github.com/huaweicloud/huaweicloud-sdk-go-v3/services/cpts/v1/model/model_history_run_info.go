package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HistoryRunInfo struct {

	// 名称
	Name *string `json:"name,omitempty"`

	// 报告id
	RunId *float64 `json:"run_id,omitempty"`

	// 任务类型（0：旧版本任务；1：融合版本任务）
	RunType *float64 `json:"run_type,omitempty"`

	// 开始时间
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间
	EndTime *string `json:"end_time,omitempty"`

	// 继续时间
	ContinueTime *float64 `json:"continue_time,omitempty"`

	// 用例或者事务名称
	TempNames *[]TempName `json:"temp_names,omitempty"`

	// 任务间用例是否并行执行
	Parallel *bool `json:"parallel,omitempty"`
}

func (o HistoryRunInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HistoryRunInfo struct{}"
	}

	return strings.Join([]string{"HistoryRunInfo", string(data)}, " ")
}
