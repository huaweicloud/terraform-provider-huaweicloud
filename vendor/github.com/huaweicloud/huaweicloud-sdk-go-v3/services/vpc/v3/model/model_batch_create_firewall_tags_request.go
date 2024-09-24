package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateFirewallTagsRequest Request Object
type BatchCreateFirewallTagsRequest struct {

	// 功能说明：ACL唯一标识 取值范围：合法UUID 约束：ID对应的ACL必须存在
	FirewallId string `json:"firewall_id"`

	Body *BatchCreateFirewallTagsRequestBody `json:"body,omitempty"`
}

func (o BatchCreateFirewallTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateFirewallTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchCreateFirewallTagsRequest", string(data)}, " ")
}
