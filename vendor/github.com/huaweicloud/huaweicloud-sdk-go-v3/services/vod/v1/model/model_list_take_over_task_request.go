package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTakeOverTaskRequest struct {

	// 任务状态。
	Status *string `json:"status,omitempty"`

	// 任务ID。
	TaskId *string `json:"task_id,omitempty"`

	// 分页编号，默认为0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。  默认10，范围[1,100]，指定task_id时该参数无效。
	Size *int32 `json:"size,omitempty"`
}

func (o ListTakeOverTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTakeOverTaskRequest struct{}"
	}

	return strings.Join([]string{"ListTakeOverTaskRequest", string(data)}, " ")
}
