package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTranscodingTaskByConsoleRequest Request Object
type DeleteTranscodingTaskByConsoleRequest struct {

	// 任务ID
	TaskId int32 `json:"task_id"`
}

func (o DeleteTranscodingTaskByConsoleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTranscodingTaskByConsoleRequest struct{}"
	}

	return strings.Join([]string{"DeleteTranscodingTaskByConsoleRequest", string(data)}, " ")
}
