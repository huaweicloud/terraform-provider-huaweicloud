package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteFirewallRequest Request Object
type DeleteFirewallRequest struct {

	// 网络ACL唯一标识
	FirewallId string `json:"firewall_id"`
}

func (o DeleteFirewallRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteFirewallRequest struct{}"
	}

	return strings.Join([]string{"DeleteFirewallRequest", string(data)}, " ")
}
