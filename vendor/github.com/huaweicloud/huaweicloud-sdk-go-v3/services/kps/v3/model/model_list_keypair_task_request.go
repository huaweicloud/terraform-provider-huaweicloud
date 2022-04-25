package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListKeypairTaskRequest struct {

	// 任务ID
	TaskId string `json:"task_id"`
}

func (o ListKeypairTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListKeypairTaskRequest struct{}"
	}

	return strings.Join([]string{"ListKeypairTaskRequest", string(data)}, " ")
}
