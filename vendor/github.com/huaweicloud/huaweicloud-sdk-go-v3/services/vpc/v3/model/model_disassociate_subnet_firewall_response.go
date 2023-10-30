package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DisassociateSubnetFirewallResponse Response Object
type DisassociateSubnetFirewallResponse struct {
	Firewall *FirewallDetail `json:"firewall,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DisassociateSubnetFirewallResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DisassociateSubnetFirewallResponse struct{}"
	}

	return strings.Join([]string{"DisassociateSubnetFirewallResponse", string(data)}, " ")
}
