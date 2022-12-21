package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DomainIpv6SwitchReq struct {

	// 域名
	Domain string `json:"domain"`

	// IPV6开关配置，默认关闭，true为开启，false为关闭
	IsIpv6 *bool `json:"is_ipv6,omitempty"`
}

func (o DomainIpv6SwitchReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainIpv6SwitchReq struct{}"
	}

	return strings.Join([]string{"DomainIpv6SwitchReq", string(data)}, " ")
}
