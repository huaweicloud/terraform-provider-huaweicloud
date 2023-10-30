package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListFirewallResponse Response Object
type ListFirewallResponse struct {

	// ACL防火墙响应体列表
	Firewalls *[]ListFirewallDetail `json:"firewalls,omitempty"`

	PageInfo *PageInfo `json:"page_info,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListFirewallResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFirewallResponse struct{}"
	}

	return strings.Join([]string{"ListFirewallResponse", string(data)}, " ")
}
