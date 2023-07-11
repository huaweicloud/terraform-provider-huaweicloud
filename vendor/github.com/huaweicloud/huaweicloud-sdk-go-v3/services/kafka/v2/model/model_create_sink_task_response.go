package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSinkTaskResponse Response Object
type CreateSinkTaskResponse struct {

	// 任务ID。
	TaskId         *string `json:"task_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateSinkTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSinkTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateSinkTaskResponse", string(data)}, " ")
}
