package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTranscodingTaskRequest struct {

	// 客户端语言
	XLanguage *string `json:"x-language,omitempty"`

	// 转码服务接受任务后产生的任务ID。一次最多10个
	TaskId *[]int64 `json:"task_id,omitempty"`

	// 任务执行状态。  取值如下： - WAITING：等待启动 - TRANSCODING：转码中 - SUCCEEDED：转码成功 - FAILED：转码失败 - CANCELED：已删除 - NEED_TO_BE_AUDIT：片源待审核
	Status *string `json:"status,omitempty"`

	// 起始时间  格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间  格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	EndTime *string `json:"end_time,omitempty"`

	// 分页编号。查询指定“task_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。查询指定“task_id”时，该参数无效。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`
}

func (o ListTranscodingTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTranscodingTaskRequest struct{}"
	}

	return strings.Join([]string{"ListTranscodingTaskRequest", string(data)}, " ")
}
