package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartTaskRequest Request Object
type StartTaskRequest struct {

	// 迁移任务ID。
	TaskId string `json:"task_id"`

	Body *StartTaskReq `json:"body,omitempty"`
}

func (o StartTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartTaskRequest struct{}"
	}

	return strings.Join([]string{"StartTaskRequest", string(data)}, " ")
}
