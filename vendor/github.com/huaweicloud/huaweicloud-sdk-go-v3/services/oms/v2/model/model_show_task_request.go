package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowTaskRequest struct {

	// 任务ID
	TaskId int64 `json:"task_id"`
}

func (o ShowTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskRequest struct{}"
	}

	return strings.Join([]string{"ShowTaskRequest", string(data)}, " ")
}
