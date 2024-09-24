package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateFirewallTagsResponse Response Object
type BatchCreateFirewallTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchCreateFirewallTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateFirewallTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchCreateFirewallTagsResponse", string(data)}, " ")
}
