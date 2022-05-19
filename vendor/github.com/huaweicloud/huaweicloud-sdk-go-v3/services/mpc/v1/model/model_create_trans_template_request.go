package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateTransTemplateRequest struct {
	Body *TransTemplate `json:"body,omitempty"`
}

func (o CreateTransTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTransTemplateRequest struct{}"
	}

	return strings.Join([]string{"CreateTransTemplateRequest", string(data)}, " ")
}
