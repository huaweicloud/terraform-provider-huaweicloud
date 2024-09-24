package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowFirewallTagsRequest Request Object
type ShowFirewallTagsRequest struct {

	// 功能说明：ACL唯一标识 取值范围：合法UUID 约束：ID对应的ACL必须存在
	FirewallId string `json:"firewall_id"`
}

func (o ShowFirewallTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowFirewallTagsRequest struct{}"
	}

	return strings.Join([]string{"ShowFirewallTagsRequest", string(data)}, " ")
}
