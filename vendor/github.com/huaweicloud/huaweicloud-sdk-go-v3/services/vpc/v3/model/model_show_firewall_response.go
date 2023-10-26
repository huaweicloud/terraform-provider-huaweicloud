package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowFirewallResponse Response Object
type ShowFirewallResponse struct {
	Firewall *FirewallDetail `json:"firewall,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowFirewallResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowFirewallResponse struct{}"
	}

	return strings.Join([]string{"ShowFirewallResponse", string(data)}, " ")
}
