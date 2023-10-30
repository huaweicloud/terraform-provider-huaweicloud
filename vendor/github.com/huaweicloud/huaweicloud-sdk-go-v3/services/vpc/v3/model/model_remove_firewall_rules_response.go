package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveFirewallRulesResponse Response Object
type RemoveFirewallRulesResponse struct {
	Firewall *FirewallDetail `json:"firewall,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o RemoveFirewallRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveFirewallRulesResponse struct{}"
	}

	return strings.Join([]string{"RemoveFirewallRulesResponse", string(data)}, " ")
}
