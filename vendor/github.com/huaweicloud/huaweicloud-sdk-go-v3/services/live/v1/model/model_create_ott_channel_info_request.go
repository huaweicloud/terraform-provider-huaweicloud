package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateOttChannelInfoRequest Request Object
type CreateOttChannelInfoRequest struct {

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-Internal访问服务。
	AccessControlAllowInternal *string `json:"Access-Control-Allow-Internal,omitempty"`

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-External访问服务。
	AccessControlAllowExternal *string `json:"Access-Control-Allow-External,omitempty"`

	Body *CreateOttChannelInfoReq `json:"body,omitempty"`
}

func (o CreateOttChannelInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOttChannelInfoRequest struct{}"
	}

	return strings.Join([]string{"CreateOttChannelInfoRequest", string(data)}, " ")
}
