package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 报告任务信息
type ReportTaskInfo struct {
	// 分钟*并发数

	Vum *float64 `json:"vum,omitempty"`
}

func (o ReportTaskInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReportTaskInfo struct{}"
	}

	return strings.Join([]string{"ReportTaskInfo", string(data)}, " ")
}
