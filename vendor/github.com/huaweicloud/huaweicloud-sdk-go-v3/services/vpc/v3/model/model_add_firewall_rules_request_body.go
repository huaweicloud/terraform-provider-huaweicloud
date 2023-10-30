package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddFirewallRulesRequestBody This is a auto create Body Object
type AddFirewallRulesRequestBody struct {
	Firewall *FirewallInsertRuleOption `json:"firewall"`
}

func (o AddFirewallRulesRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddFirewallRulesRequestBody struct{}"
	}

	return strings.Join([]string{"AddFirewallRulesRequestBody", string(data)}, " ")
}
