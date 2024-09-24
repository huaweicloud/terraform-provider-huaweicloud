package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateFirewallTagRequest Request Object
type CreateFirewallTagRequest struct {

	// 功能说明：ACL唯一标识 取值范围：合法UUID 约束：ID对应的ACL必须存在
	FirewallId string `json:"firewall_id"`

	Body *CreateFirewallTagRequestBody `json:"body,omitempty"`
}

func (o CreateFirewallTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateFirewallTagRequest struct{}"
	}

	return strings.Join([]string{"CreateFirewallTagRequest", string(data)}, " ")
}
