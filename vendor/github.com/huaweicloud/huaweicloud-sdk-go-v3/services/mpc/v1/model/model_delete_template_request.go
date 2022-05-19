package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteTemplateRequest struct {

	// 自定义转码模板ID
	TemplateId int64 `json:"template_id"`
}

func (o DeleteTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTemplateRequest struct{}"
	}

	return strings.Join([]string{"DeleteTemplateRequest", string(data)}, " ")
}
