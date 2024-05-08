package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteOttChannelInfoRequest Request Object
type DeleteOttChannelInfoRequest struct {

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-Internal访问服务。
	AccessControlAllowInternal *string `json:"Access-Control-Allow-Internal,omitempty"`

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-External访问服务。
	AccessControlAllowExternal *string `json:"Access-Control-Allow-External,omitempty"`

	// 推流域名
	Domain string `json:"domain"`

	// 组名或应用名
	AppName string `json:"app_name"`

	// 频道ID
	Id string `json:"id"`
}

func (o DeleteOttChannelInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteOttChannelInfoRequest struct{}"
	}

	return strings.Join([]string{"DeleteOttChannelInfoRequest", string(data)}, " ")
}
