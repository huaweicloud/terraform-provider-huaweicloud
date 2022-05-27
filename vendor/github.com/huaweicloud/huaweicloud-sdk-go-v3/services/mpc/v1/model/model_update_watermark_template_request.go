package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateWatermarkTemplateRequest struct {
	Body *WatermarkTemplate `json:"body,omitempty"`
}

func (o UpdateWatermarkTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateWatermarkTemplateRequest struct{}"
	}

	return strings.Join([]string{"UpdateWatermarkTemplateRequest", string(data)}, " ")
}
