package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UrlObject struct {

	// url的id
	Id *string `json:"id,omitempty"`

	// url的地址。
	Url *string `json:"url,omitempty"`

	// url的状态 processing 处理中，succeed 完成，failed 失败，waiting 等待，refreshing 刷新中，preheating 预热中。
	Status *string `json:"status,omitempty"`

	// url创建时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	CreateTime *int64 `json:"create_time,omitempty"`

	// 任务id。
	TaskId *string `json:"task_id,omitempty"`

	// 任务的类型， 其值可以为REFRESH：刷新任务、PREHEATING：预热任务、REFRESH_AFTER_PREHEATING：预热后刷新
	TaskType *string `json:"task_type,omitempty"`

	// 失败原因，url状态为failed时返回。   - ORIGIN_ERROR：源站错误。   - INNER_ERROR：内部错误。   - UNKNOWN_ERROR：未知错误。
	FailClassify *string `json:"fail_classify,omitempty"`

	// 刷新预热失败描述。
	FailDesc *string `json:"fail_desc,omitempty"`
}

func (o UrlObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UrlObject struct{}"
	}

	return strings.Join([]string{"UrlObject", string(data)}, " ")
}
