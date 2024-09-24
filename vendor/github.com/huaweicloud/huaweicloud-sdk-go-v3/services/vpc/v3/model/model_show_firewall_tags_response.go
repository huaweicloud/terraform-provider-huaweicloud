package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowFirewallTagsResponse Response Object
type ShowFirewallTagsResponse struct {

	// tag对象列表
	Tags *[]ResourceTag `json:"tags,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowFirewallTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowFirewallTagsResponse struct{}"
	}

	return strings.Join([]string{"ShowFirewallTagsResponse", string(data)}, " ")
}
