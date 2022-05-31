package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateWatermarkTemplateResponse struct {

	// 水印模板Id
	TemplateId     *int32 `json:"template_id,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateWatermarkTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateWatermarkTemplateResponse struct{}"
	}

	return strings.Join([]string{"CreateWatermarkTemplateResponse", string(data)}, " ")
}
