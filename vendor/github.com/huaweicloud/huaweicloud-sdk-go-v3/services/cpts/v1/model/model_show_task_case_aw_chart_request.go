package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTaskCaseAwChartRequest Request Object
type ShowTaskCaseAwChartRequest struct {

	// 运行任务id，即报告id。启动任务（更新任务状态或批量启停任务）接口，会返回运行任务id。
	TaskRunId int32 `json:"task_run_id"`

	// 运行用例id，通过报告管理中的“内外融合当前任务用例列表”接口获取：使用任务运行id（task_run_id）作为路径参数，可以查询到该报告关联的用例运行id集合，即返回结构体中result.case_aw_info_list[index].case_uri_i为索引为index的运行用例id（case_run_id）。
	CaseRunId int32 `json:"case_run_id"`

	// 运行用例详情id，通过报告管理中的“内外融合当前任务用例列表”接口获取：使用运行任务id（task_run_id）作为路径参数，可以查询到该报告关联的运行用例详情id集合，即返回结构体中result.case_aw_info_list[index].testcaseId为索引为index的运行用例详情id（detail_id）。
	DetailId string `json:"detail_id"`
}

func (o ShowTaskCaseAwChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskCaseAwChartRequest struct{}"
	}

	return strings.Join([]string{"ShowTaskCaseAwChartRequest", string(data)}, " ")
}
