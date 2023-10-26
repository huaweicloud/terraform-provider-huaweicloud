package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateFirewallRulesRequest Request Object
type UpdateFirewallRulesRequest struct {

	// 网络ACL唯一标识
	FirewallId string `json:"firewall_id"`

	Body *UpdateFirewallRulesRequestBody `json:"body,omitempty"`
}

func (o UpdateFirewallRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFirewallRulesRequest struct{}"
	}

	return strings.Join([]string{"UpdateFirewallRulesRequest", string(data)}, " ")
}
