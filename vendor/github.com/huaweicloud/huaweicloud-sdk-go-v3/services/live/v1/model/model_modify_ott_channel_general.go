package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelGeneral OTT频道修改通用消息
type ModifyOttChannelGeneral struct {

	// 频道推流域名
	Domain string `json:"domain"`

	// 组名或应用名
	AppName string `json:"app_name"`

	// 频道ID。频道唯一标识，为必填项
	Id string `json:"id"`

	// 频道名称
	Name string `json:"name"`
}

func (o ModifyOttChannelGeneral) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelGeneral struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelGeneral", string(data)}, " ")
}
