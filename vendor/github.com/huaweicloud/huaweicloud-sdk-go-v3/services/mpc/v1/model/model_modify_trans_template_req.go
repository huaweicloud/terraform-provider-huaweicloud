package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyTransTemplateReq struct {

	// 转码模板ID
	TemplateId int64 `json:"template_id"`

	// 转码模板名称。
	TemplateName string `json:"template_name"`

	Video *Video `json:"video,omitempty"`

	Audio *Audio `json:"audio,omitempty"`

	Common *Common `json:"common"`
}

func (o ModifyTransTemplateReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyTransTemplateReq struct{}"
	}

	return strings.Join([]string{"ModifyTransTemplateReq", string(data)}, " ")
}
