package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RemuxTask struct {

	// 任务ID
	TaskId *string `json:"task_id,omitempty"`

	// 任务状态。  取值如下： - INIT：初始状态。 - WAITING：等待启动。 - PROCESSING：处理中。 - SUCCEED：处理成功。 - FAILED：处理失败。 - CANCELED：已取消。
	Status *string `json:"status,omitempty"`

	// 任务创建时间
	CreateTime *string `json:"create_time,omitempty"`

	// 任务启动时间
	StartTime *string `json:"start_time,omitempty"`

	// 任务结束时间
	EndTime *string `json:"end_time,omitempty"`

	// 任务的返回码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 错误描述
	Description *string `json:"description,omitempty"`

	// 用户数据。
	UserData *string `json:"user_data,omitempty"`

	Input *ObsObjInfo `json:"input,omitempty"`

	Output *ObsObjInfo `json:"output,omitempty"`

	OutputParam *RemuxOutputParam `json:"output_param,omitempty"`

	// 任务完成进度百分比值。
	CompleteRatio *int32 `json:"complete_ratio,omitempty"`

	OutputMetadata *MetaData `json:"output_metadata,omitempty"`
}

func (o RemuxTask) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemuxTask struct{}"
	}

	return strings.Join([]string{"RemuxTask", string(data)}, " ")
}
