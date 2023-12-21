package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// LoginType 登录类型，包含如下: - \"mysql\" # mysql服务 - \"rdp\" # rdp服务服务 - \"ssh\" # ssh服务 - \"vsftp\" # vsftp服务
type LoginType struct {
}

func (o LoginType) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LoginType struct{}"
	}

	return strings.Join([]string{"LoginType", string(data)}, " ")
}
