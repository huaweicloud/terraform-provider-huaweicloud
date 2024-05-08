package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BlockedIpRequestInfo 解除拦截的IP详情
type BlockedIpRequestInfo struct {

	// 主机ID
	HostId string `json:"host_id"`

	// 攻击源IP
	SrcIp string `json:"src_ip"`

	// 登录类型，包含如下: - \"mysql\" # mysql服务 - \"rdp\" # rdp服务服务 - \"ssh\" # ssh服务 - \"vsftp\" # vsftp服务
	LoginType string `json:"login_type"`
}

func (o BlockedIpRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BlockedIpRequestInfo struct{}"
	}

	return strings.Join([]string{"BlockedIpRequestInfo", string(data)}, " ")
}
