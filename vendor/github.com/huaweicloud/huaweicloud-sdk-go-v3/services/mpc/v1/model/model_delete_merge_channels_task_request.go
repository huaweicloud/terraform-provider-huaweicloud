package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteMergeChannelsTaskRequest struct {

	// 任务ID
	TaskId string `json:"task_id"`
}

func (o DeleteMergeChannelsTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteMergeChannelsTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteMergeChannelsTaskRequest", string(data)}, " ")
}
