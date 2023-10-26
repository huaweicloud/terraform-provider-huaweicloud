package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DisassociateSubnetFirewallRequestBody This is a auto create Body Object
type DisassociateSubnetFirewallRequestBody struct {

	// 解绑ACL的子网列表
	Subnets []FirewallAssociation `json:"subnets"`
}

func (o DisassociateSubnetFirewallRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DisassociateSubnetFirewallRequestBody struct{}"
	}

	return strings.Join([]string{"DisassociateSubnetFirewallRequestBody", string(data)}, " ")
}
