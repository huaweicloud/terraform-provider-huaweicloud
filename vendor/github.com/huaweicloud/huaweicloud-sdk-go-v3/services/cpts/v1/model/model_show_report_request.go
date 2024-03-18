package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowReportRequest Request Object
type ShowReportRequest struct {

	// 运行任务id，即报告id。启动任务（更新任务状态或批量启停任务）接口，会返回运行任务id。
	TaskRunId int32 `json:"task_run_id"`

	// 运行用例id，报告管理中的“内外融合当前任务用例列表”接口，使用任务运行id（task_run_id）作为路径参数，可以查询到该报告关联的用例运行id集合，即返回结构体中result.case_aw_info_list[index].case_uri_i为索引为index的运行用例id（case_run_id）。
	CaseRunId int32 `json:"case_run_id"`

	// 曲线图点数
	BrokensLimitCount int32 `json:"brokens_limit_count"`
}

func (o ShowReportRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowReportRequest struct{}"
	}

	return strings.Join([]string{"ShowReportRequest", string(data)}, " ")
}
