package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateTranscodingTaskRequest struct {
	Body *CreateTranscodingReq `json:"body,omitempty"`
}

func (o CreateTranscodingTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTranscodingTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateTranscodingTaskRequest", string(data)}, " ")
}
