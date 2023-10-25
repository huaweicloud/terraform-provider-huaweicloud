package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddFirewallRulesResponse Response Object
type AddFirewallRulesResponse struct {
	Firewall *FirewallDetail `json:"firewall,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AddFirewallRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddFirewallRulesResponse struct{}"
	}

	return strings.Join([]string{"AddFirewallRulesResponse", string(data)}, " ")
}
