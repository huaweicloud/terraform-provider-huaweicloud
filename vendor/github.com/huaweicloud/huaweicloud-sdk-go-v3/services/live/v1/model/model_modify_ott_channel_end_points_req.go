package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelEndPointsReq OTT频道修改打包消息体
type ModifyOttChannelEndPointsReq struct {

	// 频道推流域名
	Domain string `json:"domain"`

	// 组名或应用名
	AppName string `json:"app_name"`

	// 频道ID。频道唯一标识，为必填项
	Id string `json:"id"`

	// 频道出流信息
	Endpoints *[]EndpointItem `json:"endpoints,omitempty"`
}

func (o ModifyOttChannelEndPointsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelEndPointsReq struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelEndPointsReq", string(data)}, " ")
}
