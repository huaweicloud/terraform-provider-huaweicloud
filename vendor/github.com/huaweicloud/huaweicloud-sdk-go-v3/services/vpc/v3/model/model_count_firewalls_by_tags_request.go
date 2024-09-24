package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CountFirewallsByTagsRequest Request Object
type CountFirewallsByTagsRequest struct {
	Body *CountFirewallsByTagsRequestBody `json:"body,omitempty"`
}

func (o CountFirewallsByTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CountFirewallsByTagsRequest struct{}"
	}

	return strings.Join([]string{"CountFirewallsByTagsRequest", string(data)}, " ")
}
