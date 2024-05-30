package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowUrlTaskInfoRequest Request Object
type ShowUrlTaskInfoRequest struct {

	// 起始时间戳（毫秒），默认当天00:00。
	StartTime *int64 `json:"start_time,omitempty"`

	// 结束时间戳（毫秒），默认次日00:00。
	EndTime *int64 `json:"end_time,omitempty"`

	// 偏移量：特定数据字段与起始数据字段位置的距离，默认为0。
	Offset *int32 `json:"offset,omitempty"`

	// 单次查询数据条数，上限为100，默认为10。
	Limit *int32 `json:"limit,omitempty"`

	// 刷新预热url。
	Url *string `json:"url,omitempty"`

	// 任务类型，REFRESH：刷新任务；PREHEATING：预热任务。
	TaskType *string `json:"task_type,omitempty"`

	// url状态，状态类型：processing：处理中；succeed：完成；failed：失败；waiting：等待；refreshing：刷新中; preheating : 预热中。
	Status *string `json:"status,omitempty"`

	// 文件类型，file:文件;directory:目录。
	FileType *string `json:"file_type,omitempty"`
}

func (o ShowUrlTaskInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUrlTaskInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowUrlTaskInfoRequest", string(data)}, " ")
}
