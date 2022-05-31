package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowTakeOverTaskDetailsRequest struct {

	// 任务ID。
	TaskId string `json:"task_id"`

	// 分页编号，默认为0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。  默认10，范围[1,100]
	Size *int32 `json:"size,omitempty"`
}

func (o ShowTakeOverTaskDetailsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTakeOverTaskDetailsRequest struct{}"
	}

	return strings.Join([]string{"ShowTakeOverTaskDetailsRequest", string(data)}, " ")
}
