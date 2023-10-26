package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateFirewallRequestBody This is a auto create Body Object
type UpdateFirewallRequestBody struct {
	Firewall *UpdateFirewallOption `json:"firewall"`
}

func (o UpdateFirewallRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFirewallRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateFirewallRequestBody", string(data)}, " ")
}
