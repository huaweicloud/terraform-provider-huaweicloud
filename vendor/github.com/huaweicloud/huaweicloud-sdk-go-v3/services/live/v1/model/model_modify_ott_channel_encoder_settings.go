package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelEncoderSettings OTT频道通用消息
type ModifyOttChannelEncoderSettings struct {

	// 频道推流域名
	Domain string `json:"domain"`

	// 组名或应用名
	AppName string `json:"app_name"`

	// 频道ID。频道唯一标识，为必填项
	Id string `json:"id"`

	// 转码模板配置
	EncoderSettings *[]ModifyOttChannelEncoderSettingsEncoderSettings `json:"encoder_settings,omitempty"`
}

func (o ModifyOttChannelEncoderSettings) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelEncoderSettings struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelEncoderSettings", string(data)}, " ")
}
