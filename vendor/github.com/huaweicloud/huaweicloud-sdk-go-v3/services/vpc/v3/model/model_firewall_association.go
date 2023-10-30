package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FirewallAssociation
type FirewallAssociation struct {

	// 功能说明：ACL绑定的子网ID
	VirsubnetId string `json:"virsubnet_id"`
}

func (o FirewallAssociation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FirewallAssociation struct{}"
	}

	return strings.Join([]string{"FirewallAssociation", string(data)}, " ")
}
