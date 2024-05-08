package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelInfoEncoderSettingsRequest Request Object
type ModifyOttChannelInfoEncoderSettingsRequest struct {

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-Internal访问服务。
	AccessControlAllowInternal *string `json:"Access-Control-Allow-Internal,omitempty"`

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-External访问服务。
	AccessControlAllowExternal *string `json:"Access-Control-Allow-External,omitempty"`

	Body *ModifyOttChannelEncoderSettings `json:"body,omitempty"`
}

func (o ModifyOttChannelInfoEncoderSettingsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelInfoEncoderSettingsRequest struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelInfoEncoderSettingsRequest", string(data)}, " ")
}
