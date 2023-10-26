package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FirewallRemoveRuleItemOption
type FirewallRemoveRuleItemOption struct {

	// 功能说明：要删除的ACL规则id
	Id string `json:"id"`
}

func (o FirewallRemoveRuleItemOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FirewallRemoveRuleItemOption struct{}"
	}

	return strings.Join([]string{"FirewallRemoveRuleItemOption", string(data)}, " ")
}
