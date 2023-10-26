package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FirewallInsertRuleOption
type FirewallInsertRuleOption struct {

	// 功能说明：ACL添加入方向规则列表
	IngressRules *[]FirewallInsertRuleItemOption `json:"ingress_rules,omitempty"`

	// 功能说明：ACL添加出方向规则列表
	EgressRules *[]FirewallInsertRuleItemOption `json:"egress_rules,omitempty"`

	// 功能说明：插入ACL的规则在入方向或者出方向某条规则位置后，不指定则在入方向或者出方向规则列表最前面插入规则 约束：指定了insert_after_rule，ingress_rules和egress_rules只能同时设置一个，且该规则在入方向或者出方向规则中存在
	InsertAfterRule *string `json:"insert_after_rule,omitempty"`
}

func (o FirewallInsertRuleOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FirewallInsertRuleOption struct{}"
	}

	return strings.Join([]string{"FirewallInsertRuleOption", string(data)}, " ")
}
