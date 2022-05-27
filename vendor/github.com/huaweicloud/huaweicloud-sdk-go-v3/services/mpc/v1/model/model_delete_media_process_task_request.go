package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteMediaProcessTaskRequest struct {

	// 任务ID
	TaskId string `json:"task_id"`
}

func (o DeleteMediaProcessTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMediaProcessTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteMediaProcessTaskRequest", string(data)}, " ")
}
