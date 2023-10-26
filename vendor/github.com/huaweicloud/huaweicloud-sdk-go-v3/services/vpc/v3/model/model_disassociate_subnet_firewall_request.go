package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DisassociateSubnetFirewallRequest Request Object
type DisassociateSubnetFirewallRequest struct {

	// 网络ACL唯一标识
	FirewallId string `json:"firewall_id"`

	Body *DisassociateSubnetFirewallRequestBody `json:"body,omitempty"`
}

func (o DisassociateSubnetFirewallRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DisassociateSubnetFirewallRequest struct{}"
	}

	return strings.Join([]string{"DisassociateSubnetFirewallRequest", string(data)}, " ")
}
