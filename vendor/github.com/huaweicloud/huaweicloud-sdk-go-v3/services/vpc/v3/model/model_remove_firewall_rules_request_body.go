package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveFirewallRulesRequestBody This is a auto create Body Object
type RemoveFirewallRulesRequestBody struct {
	Firewall *FirewallRemoveRuleOption `json:"firewall"`
}

func (o RemoveFirewallRulesRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveFirewallRulesRequestBody struct{}"
	}

	return strings.Join([]string{"RemoveFirewallRulesRequestBody", string(data)}, " ")
}
