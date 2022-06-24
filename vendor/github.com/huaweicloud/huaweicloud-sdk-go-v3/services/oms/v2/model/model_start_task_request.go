package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type StartTaskRequest struct {

	// 迁移任务ID。
	TaskId int64 `json:"task_id"`

	Body *StartTaskReq `json:"body,omitempty"`
}

func (o StartTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartTaskRequest struct{}"
	}

	return strings.Join([]string{"StartTaskRequest", string(data)}, " ")
}
