package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateFirewallTagResponse Response Object
type CreateFirewallTagResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreateFirewallTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateFirewallTagResponse struct{}"
	}

	return strings.Join([]string{"CreateFirewallTagResponse", string(data)}, " ")
}
