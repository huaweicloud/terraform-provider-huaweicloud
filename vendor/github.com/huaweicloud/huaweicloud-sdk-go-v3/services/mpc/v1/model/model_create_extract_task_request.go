package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateExtractTaskRequest struct {
	Body *CreateExtractTaskReq `json:"body,omitempty"`
}

func (o CreateExtractTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateExtractTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateExtractTaskRequest", string(data)}, " ")
}
