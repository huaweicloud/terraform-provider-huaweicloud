package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTranscodingTaskResponse struct {

	// 任务ID。 如果返回值为200 OK，为接受任务后产生的任务ID。
	TaskId         *int32 `json:"task_id,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateTranscodingTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTranscodingTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateTranscodingTaskResponse", string(data)}, " ")
}
