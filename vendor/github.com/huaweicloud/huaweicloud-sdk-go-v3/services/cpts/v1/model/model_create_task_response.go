package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTaskResponse struct {
	// code

	Code *string `json:"code,omitempty"`
	// message

	Message *string `json:"message,omitempty"`
	// task_id

	TaskId         *int32 `json:"task_id,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateTaskResponse", string(data)}, " ")
}
