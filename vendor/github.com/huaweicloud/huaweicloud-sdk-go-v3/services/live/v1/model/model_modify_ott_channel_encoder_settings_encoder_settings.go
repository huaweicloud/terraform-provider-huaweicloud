package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelEncoderSettingsEncoderSettings 转码模板配置项
type ModifyOttChannelEncoderSettingsEncoderSettings struct {

	// 转码模板ID
	TemplateId *string `json:"template_id,omitempty"`
}

func (o ModifyOttChannelEncoderSettingsEncoderSettings) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelEncoderSettingsEncoderSettings struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelEncoderSettingsEncoderSettings", string(data)}, " ")
}
