package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CommonCreateTaskRsp struct {

	// 任务ID
	TaskId *string `json:"task_id,omitempty"`
}

func (o CommonCreateTaskRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CommonCreateTaskRsp struct{}"
	}

	return strings.Join([]string{"CommonCreateTaskRsp", string(data)}, " ")
}
