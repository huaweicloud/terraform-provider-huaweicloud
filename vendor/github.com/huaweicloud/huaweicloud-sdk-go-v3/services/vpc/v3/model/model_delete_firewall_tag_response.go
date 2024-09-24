package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteFirewallTagResponse Response Object
type DeleteFirewallTagResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteFirewallTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteFirewallTagResponse struct{}"
	}

	return strings.Join([]string{"DeleteFirewallTagResponse", string(data)}, " ")
}
