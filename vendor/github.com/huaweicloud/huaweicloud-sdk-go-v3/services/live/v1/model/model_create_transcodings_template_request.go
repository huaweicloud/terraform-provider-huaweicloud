package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateTranscodingsTemplateRequest struct {
	Body *StreamTranscodingTemplate `json:"body,omitempty"`
}

func (o CreateTranscodingsTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTranscodingsTemplateRequest struct{}"
	}

	return strings.Join([]string{"CreateTranscodingsTemplateRequest", string(data)}, " ")
}
