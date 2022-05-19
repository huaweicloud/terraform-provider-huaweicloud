package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteExtractTaskRequest struct {

	// 任务ID
	TaskId string `json:"task_id"`
}

func (o DeleteExtractTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteExtractTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteExtractTaskRequest", string(data)}, " ")
}
