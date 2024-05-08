package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelInfoGeneralRequest Request Object
type ModifyOttChannelInfoGeneralRequest struct {

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-Internal访问服务。
	AccessControlAllowInternal *string `json:"Access-Control-Allow-Internal,omitempty"`

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-External访问服务。
	AccessControlAllowExternal *string `json:"Access-Control-Allow-External,omitempty"`

	Body *ModifyOttChannelGeneral `json:"body,omitempty"`
}

func (o ModifyOttChannelInfoGeneralRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelInfoGeneralRequest struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelInfoGeneralRequest", string(data)}, " ")
}
