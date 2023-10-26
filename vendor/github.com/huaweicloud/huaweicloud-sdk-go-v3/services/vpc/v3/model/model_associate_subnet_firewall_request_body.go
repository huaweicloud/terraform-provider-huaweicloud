package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssociateSubnetFirewallRequestBody This is a auto create Body Object
type AssociateSubnetFirewallRequestBody struct {

	// 绑定ACL的子网列表
	Subnets []FirewallAssociation `json:"subnets"`
}

func (o AssociateSubnetFirewallRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateSubnetFirewallRequestBody struct{}"
	}

	return strings.Join([]string{"AssociateSubnetFirewallRequestBody", string(data)}, " ")
}
