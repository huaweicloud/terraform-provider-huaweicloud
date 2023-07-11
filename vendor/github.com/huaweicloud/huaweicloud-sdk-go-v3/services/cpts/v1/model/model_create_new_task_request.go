package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateNewTaskRequest Request Object
type CreateNewTaskRequest struct {
	Body *NewTaskInfo `json:"body,omitempty"`
}

func (o CreateNewTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateNewTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateNewTaskRequest", string(data)}, " ")
}
