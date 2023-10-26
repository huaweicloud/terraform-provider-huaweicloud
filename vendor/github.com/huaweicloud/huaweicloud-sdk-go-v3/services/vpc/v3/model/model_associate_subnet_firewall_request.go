package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssociateSubnetFirewallRequest Request Object
type AssociateSubnetFirewallRequest struct {

	// 网络ACL唯一标识
	FirewallId string `json:"firewall_id"`

	Body *AssociateSubnetFirewallRequestBody `json:"body,omitempty"`
}

func (o AssociateSubnetFirewallRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateSubnetFirewallRequest struct{}"
	}

	return strings.Join([]string{"AssociateSubnetFirewallRequest", string(data)}, " ")
}
