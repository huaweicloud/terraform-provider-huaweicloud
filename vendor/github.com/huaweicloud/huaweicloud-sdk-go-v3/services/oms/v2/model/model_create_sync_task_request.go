package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSyncTaskRequest Request Object
type CreateSyncTaskRequest struct {
	Body *CreateSyncTaskReq `json:"body,omitempty"`
}

func (o CreateSyncTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSyncTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateSyncTaskRequest", string(data)}, " ")
}
