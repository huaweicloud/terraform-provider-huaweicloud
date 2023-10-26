package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FirewallUpdateRuleOption
type FirewallUpdateRuleOption struct {

	// 功能说明：ACL更新入方向规则列表 约束：ingress_rules和egress_rules仅能同时设置一个，且当前只支持同时更新一条规则
	IngressRules *[]FirewallUpdateRuleItemOption `json:"ingress_rules,omitempty"`

	// 功能说明：ACL更新出方向规则列表 约束：ingress_rules和egress_rules仅能同时设置一个，且当前只支持同时更新一条规则
	EgressRules *[]FirewallUpdateRuleItemOption `json:"egress_rules,omitempty"`
}

func (o FirewallUpdateRuleOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FirewallUpdateRuleOption struct{}"
	}

	return strings.Join([]string{"FirewallUpdateRuleOption", string(data)}, " ")
}
