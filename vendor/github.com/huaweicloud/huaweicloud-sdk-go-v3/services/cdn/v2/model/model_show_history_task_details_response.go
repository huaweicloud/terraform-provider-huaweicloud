package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowHistoryTaskDetailsResponse Response Object
type ShowHistoryTaskDetailsResponse struct {

	// 任务id。
	Id *string `json:"id,omitempty"`

	// 任务类型，refresh：刷新任务；preheating：预热任务。
	TaskType *string `json:"task_type,omitempty"`

	// 任务执行结果,task_done:成功，task_inprocess:处理中。
	Status *string `json:"status,omitempty"`

	// 本次提交的url列表。
	Urls *[]UrlObject `json:"urls,omitempty"`

	// 创建时间。
	CreateTime *int64 `json:"create_time,omitempty"`

	// 处理中的url个数。
	Processing *int32 `json:"processing,omitempty"`

	// 成功处理的url个数。
	Succeed *int32 `json:"succeed,omitempty"`

	// 处理失败的url个数。
	Failed *int32 `json:"failed,omitempty"`

	// 历史任务的url个数。
	Total *int32 `json:"total,omitempty"`

	// 文件类型，file：文件；directory：目录，默认是文件file。
	FileType *string `json:"file_type,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowHistoryTaskDetailsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowHistoryTaskDetailsResponse struct{}"
	}

	return strings.Join([]string{"ShowHistoryTaskDetailsResponse", string(data)}, " ")
}
