package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssociateSubnetFirewallResponse Response Object
type AssociateSubnetFirewallResponse struct {
	Firewall *FirewallDetail `json:"firewall,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AssociateSubnetFirewallResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateSubnetFirewallResponse struct{}"
	}

	return strings.Join([]string{"AssociateSubnetFirewallResponse", string(data)}, " ")
}
