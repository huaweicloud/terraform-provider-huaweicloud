package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TemplateInfo struct {

	// 转码模板ID。
	TemplateId *int32 `json:"template_id,omitempty"`

	Template *QueryTransTemplate `json:"template,omitempty"`
}

func (o TemplateInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TemplateInfo struct{}"
	}

	return strings.Join([]string{"TemplateInfo", string(data)}, " ")
}
