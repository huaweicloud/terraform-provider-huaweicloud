package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 模板响应
type TemplateResponse struct {
	Template *TemplateResponseBody `json:"template,omitempty"`
}

func (o TemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TemplateResponse struct{}"
	}

	return strings.Join([]string{"TemplateResponse", string(data)}, " ")
}
