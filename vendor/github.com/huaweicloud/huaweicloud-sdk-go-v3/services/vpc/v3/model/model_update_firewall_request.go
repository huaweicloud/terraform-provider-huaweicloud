package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateFirewallRequest Request Object
type UpdateFirewallRequest struct {

	// 网络ACL的唯一标识
	FirewallId string `json:"firewall_id"`

	Body *UpdateFirewallRequestBody `json:"body,omitempty"`
}

func (o UpdateFirewallRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFirewallRequest struct{}"
	}

	return strings.Join([]string{"UpdateFirewallRequest", string(data)}, " ")
}
