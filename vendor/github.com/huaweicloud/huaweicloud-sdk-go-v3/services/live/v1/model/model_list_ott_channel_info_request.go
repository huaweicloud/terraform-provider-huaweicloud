package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOttChannelInfoRequest Request Object
type ListOttChannelInfoRequest struct {

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-Internal访问服务。
	AccessControlAllowInternal *string `json:"Access-Control-Allow-Internal,omitempty"`

	// 服务鉴权Token，服务开启鉴权，必须携带Access-Control-Allow-External访问服务。
	AccessControlAllowExternal *string `json:"Access-Control-Allow-External,omitempty"`

	// 推流域名
	Domain *string `json:"domain,omitempty"`

	// 组名或应用名
	AppName *string `json:"app_name,omitempty"`

	// 频道ID
	Id *string `json:"id,omitempty"`

	// 每页记录数，取值范围[1,100]，默认值10
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量。表示从此偏移量开始查询，offset大于等于0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListOttChannelInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOttChannelInfoRequest struct{}"
	}

	return strings.Join([]string{"ListOttChannelInfoRequest", string(data)}, " ")
}
