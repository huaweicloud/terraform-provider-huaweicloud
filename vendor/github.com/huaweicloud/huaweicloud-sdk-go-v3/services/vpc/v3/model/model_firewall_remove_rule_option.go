package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FirewallRemoveRuleOption
type FirewallRemoveRuleOption struct {

	// 功能说明：ACL删除入方向规则列表 约束：ingress_rules和egress_rules仅能同时设置一个
	IngressRules *[]FirewallRemoveRuleItemOption `json:"ingress_rules,omitempty"`

	// 功能说明：ACL删除出方向规则列表 约束：ingress_rules和egress_rules仅能同时设置一个
	EgressRules *[]FirewallRemoveRuleItemOption `json:"egress_rules,omitempty"`
}

func (o FirewallRemoveRuleOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FirewallRemoveRuleOption struct{}"
	}

	return strings.Join([]string{"FirewallRemoveRuleOption", string(data)}, " ")
}
