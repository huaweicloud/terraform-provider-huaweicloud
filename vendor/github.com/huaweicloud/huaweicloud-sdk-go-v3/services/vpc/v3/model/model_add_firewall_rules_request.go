package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddFirewallRulesRequest Request Object
type AddFirewallRulesRequest struct {

	// 网络ACL的唯一标识
	FirewallId string `json:"firewall_id"`

	Body *AddFirewallRulesRequestBody `json:"body,omitempty"`
}

func (o AddFirewallRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddFirewallRulesRequest struct{}"
	}

	return strings.Join([]string{"AddFirewallRulesRequest", string(data)}, " ")
}
