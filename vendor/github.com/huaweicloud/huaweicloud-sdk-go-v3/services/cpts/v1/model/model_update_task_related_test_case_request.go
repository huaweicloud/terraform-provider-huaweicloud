package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTaskRelatedTestCaseRequest Request Object
type UpdateTaskRelatedTestCaseRequest struct {

	// 任务id
	TaskId int32 `json:"task_id"`

	Body *UpdateNewTaskRequestBody `json:"body,omitempty"`
}

func (o UpdateTaskRelatedTestCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskRelatedTestCaseRequest struct{}"
	}

	return strings.Join([]string{"UpdateTaskRelatedTestCaseRequest", string(data)}, " ")
}
