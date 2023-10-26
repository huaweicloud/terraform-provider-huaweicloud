package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateFirewallRulesRequestBody This is a auto create Body Object
type UpdateFirewallRulesRequestBody struct {
	Firewall *FirewallUpdateRuleOption `json:"firewall"`
}

func (o UpdateFirewallRulesRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFirewallRulesRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateFirewallRulesRequestBody", string(data)}, " ")
}
