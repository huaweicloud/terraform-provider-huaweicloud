package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateWatermarkTemplateRequest struct {
	Body *CreateWatermarkTemplateReq `json:"body,omitempty"`
}

func (o CreateWatermarkTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateWatermarkTemplateRequest struct{}"
	}

	return strings.Join([]string{"CreateWatermarkTemplateRequest", string(data)}, " ")
}
