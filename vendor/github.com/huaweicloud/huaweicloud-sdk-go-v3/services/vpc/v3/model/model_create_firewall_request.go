package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateFirewallRequest Request Object
type CreateFirewallRequest struct {
	Body *CreateFirewallRequestBody `json:"body,omitempty"`
}

func (o CreateFirewallRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateFirewallRequest struct{}"
	}

	return strings.Join([]string{"CreateFirewallRequest", string(data)}, " ")
}
