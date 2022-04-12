package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateCaseRequest struct {
	Body *CreateCaseRequestBody `json:"body,omitempty"`
}

func (o CreateCaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCaseRequest struct{}"
	}

	return strings.Join([]string{"CreateCaseRequest", string(data)}, " ")
}
