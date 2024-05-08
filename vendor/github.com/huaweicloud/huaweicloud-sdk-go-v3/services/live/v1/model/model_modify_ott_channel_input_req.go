package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelInputReq OTT频道修改入流消息体
type ModifyOttChannelInputReq struct {

	// 频道推流域名
	Domain string `json:"domain"`

	// 组名或应用名
	AppName string `json:"app_name"`

	// 频道ID。频道唯一标识，为必填项
	Id string `json:"id"`

	Input *InputStreamInfo `json:"input,omitempty"`
}

func (o ModifyOttChannelInputReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelInputReq struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelInputReq", string(data)}, " ")
}
