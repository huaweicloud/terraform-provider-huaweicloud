package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateFirewallRulesResponse Response Object
type UpdateFirewallRulesResponse struct {
	Firewall *FirewallDetail `json:"firewall,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateFirewallRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFirewallRulesResponse struct{}"
	}

	return strings.Join([]string{"UpdateFirewallRulesResponse", string(data)}, " ")
}
