package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchUpdateTasksRequest Request Object
type BatchUpdateTasksRequest struct {
	Body *BatchUpdateTasksReq `json:"body,omitempty"`
}

func (o BatchUpdateTasksRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchUpdateTasksRequest struct{}"
	}

	return strings.Join([]string{"BatchUpdateTasksRequest", string(data)}, " ")
}
