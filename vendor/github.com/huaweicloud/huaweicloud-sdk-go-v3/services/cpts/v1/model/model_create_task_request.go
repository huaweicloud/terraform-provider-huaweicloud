package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateTaskRequest struct {
	Body *CreateTaskRequestBody `json:"body,omitempty"`
}

func (o CreateTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateTaskRequest", string(data)}, " ")
}
