package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MbTasksReportReq struct {

	// 任务ID。 如果返回值为200 OK，为接受任务后产生的任务ID。
	TaskId *string `json:"task_id,omitempty"`

	// 任务执行状态。 取值为RUNNING/FINISHED/FAILED。
	Status *string `json:"status,omitempty"`

	// 任务名称。 取值为RESET_TRACKS/MERGE_CHANNELS。
	TaskName *string `json:"task_name,omitempty"`

	// 失败任务是否重试。
	Retry *bool `json:"retry,omitempty"`

	Parameter *MbTaskParameter `json:"parameter,omitempty"`
}

func (o MbTasksReportReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MbTasksReportReq struct{}"
	}

	return strings.Join([]string{"MbTasksReportReq", string(data)}, " ")
}
