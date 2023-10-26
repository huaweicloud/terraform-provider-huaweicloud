package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateFirewallResponse Response Object
type CreateFirewallResponse struct {
	Firewall *FirewallDetail `json:"firewall,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateFirewallResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateFirewallResponse struct{}"
	}

	return strings.Join([]string{"CreateFirewallResponse", string(data)}, " ")
}
