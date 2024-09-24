package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteFirewallTagRequest Request Object
type DeleteFirewallTagRequest struct {

	// 功能说明：ACL唯一标识 取值范围：合法UUID 约束：ID对应的ACL必须存在
	FirewallId string `json:"firewall_id"`

	// 功能说明：标签键
	TagKey string `json:"tag_key"`
}

func (o DeleteFirewallTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteFirewallTagRequest struct{}"
	}

	return strings.Join([]string{"DeleteFirewallTagRequest", string(data)}, " ")
}
