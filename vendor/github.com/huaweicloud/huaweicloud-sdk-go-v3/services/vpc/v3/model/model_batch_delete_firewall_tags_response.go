package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteFirewallTagsResponse Response Object
type BatchDeleteFirewallTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchDeleteFirewallTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteFirewallTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteFirewallTagsResponse", string(data)}, " ")
}
