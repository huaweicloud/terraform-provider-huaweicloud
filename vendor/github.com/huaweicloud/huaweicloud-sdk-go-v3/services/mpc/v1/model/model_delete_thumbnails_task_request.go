package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteThumbnailsTaskRequest struct {

	// 任务ID
	TaskId string `json:"task_id"`
}

func (o DeleteThumbnailsTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteThumbnailsTaskRequest struct{}"
	}

	return strings.Join([]string{"DeleteThumbnailsTaskRequest", string(data)}, " ")
}
