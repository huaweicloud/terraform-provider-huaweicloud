package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowFirewallRequest Request Object
type ShowFirewallRequest struct {

	// 网络ACL的唯一标识
	FirewallId string `json:"firewall_id"`
}

func (o ShowFirewallRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowFirewallRequest struct{}"
	}

	return strings.Join([]string{"ShowFirewallRequest", string(data)}, " ")
}
