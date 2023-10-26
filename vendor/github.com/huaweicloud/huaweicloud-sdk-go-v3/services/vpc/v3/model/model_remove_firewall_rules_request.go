package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveFirewallRulesRequest Request Object
type RemoveFirewallRulesRequest struct {

	// 网络ACL唯一标识
	FirewallId string `json:"firewall_id"`

	Body *RemoveFirewallRulesRequestBody `json:"body,omitempty"`
}

func (o RemoveFirewallRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveFirewallRulesRequest struct{}"
	}

	return strings.Join([]string{"RemoveFirewallRulesRequest", string(data)}, " ")
}
