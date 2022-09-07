package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateTranscodeTemplateRequest struct {
	Body *ModifyTransTemplate `json:"body,omitempty"`
}

func (o UpdateTranscodeTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTranscodeTemplateRequest struct{}"
	}

	return strings.Join([]string{"UpdateTranscodeTemplateRequest", string(data)}, " ")
}
