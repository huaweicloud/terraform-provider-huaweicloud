package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RelatedTempRunningData struct {

	// 运行任务id，即报告id。启动任务（更新任务状态或批量启停任务）接口，会返回运行任务id。
	TaskRunInfoId *int32 `json:"task_run_info_id,omitempty"`

	// 运行用例id。对应其他（如报告）接口的运行用例id（case_run_id）。
	RelatedTempRunningId *int32 `json:"related_temp_running_id,omitempty"`

	// 用例id
	TempId *int32 `json:"temp_id,omitempty"`

	// 用例名称
	TempName *string `json:"temp_name,omitempty"`

	// 请求信息，包括请求名称，方法，url信息
	ContentMethodUrl *[]string `json:"content_method_url,omitempty"`

	// 最近一次运行的报告简略信息
	RelatedTempRunningData *[]TempRunningData `json:"related_temp_running_data,omitempty"`
}

func (o RelatedTempRunningData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RelatedTempRunningData struct{}"
	}

	return strings.Join([]string{"RelatedTempRunningData", string(data)}, " ")
}
