package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TempRunningData struct {

	// 请求信息，包括请求名称，方法，url信息
	ContentMethodUrl *[]string `json:"content_method_url,omitempty"`

	// 请求运行状态（0：正常返回；1：解析失败； 2：比对失败； 3：响应超时；）
	CrawlerStatus *int32 `json:"crawler_status,omitempty"`

	// 运行用例id。对应其他（如报告）接口的运行用例id（case_run_id）。
	RelatedTempRunningId *int32 `json:"related_temp_running_id,omitempty"`

	// 运行任务id，即报告id。启动任务（更新任务状态或批量启停任务）接口，会返回运行任务id。
	TaskRunInfoId *int32 `json:"task_run_info_id,omitempty"`

	// 用例或者事务id
	TempId *int32 `json:"temp_id,omitempty"`

	// 用例或者事务名称
	TempName *string `json:"temp_name,omitempty"`

	// 运行状态（9：表示等待运行；0：表示运行中；2：表示结束；3：异常中止；4：用户主动终止（完成状态）；5：用户主动终止（终止中状态））
	TempRunningStatus *int32 `json:"temp_running_status,omitempty"`
}

func (o TempRunningData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TempRunningData struct{}"
	}

	return strings.Join([]string{"TempRunningData", string(data)}, " ")
}
