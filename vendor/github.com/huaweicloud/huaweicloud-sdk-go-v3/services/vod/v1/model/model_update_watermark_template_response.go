package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateWatermarkTemplateResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateWatermarkTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateWatermarkTemplateResponse struct{}"
	}

	return strings.Join([]string{"UpdateWatermarkTemplateResponse", string(data)}, " ")
}
