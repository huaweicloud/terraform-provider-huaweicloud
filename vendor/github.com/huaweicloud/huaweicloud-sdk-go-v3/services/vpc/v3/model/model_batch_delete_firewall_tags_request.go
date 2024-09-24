package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteFirewallTagsRequest Request Object
type BatchDeleteFirewallTagsRequest struct {

	// 功能说明：ACL唯一标识 取值范围：合法UUID 约束：ID对应的ACL必须存在
	FirewallId string `json:"firewall_id"`

	Body *BatchDeleteFirewallTagsRequestBody `json:"body,omitempty"`
}

func (o BatchDeleteFirewallTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteFirewallTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchDeleteFirewallTagsRequest", string(data)}, " ")
}
