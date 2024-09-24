package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateFirewallTagRequestBody This is a auto create Body Object
type CreateFirewallTagRequestBody struct {
	Tag *ResourceTag `json:"tag"`
}

func (o CreateFirewallTagRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateFirewallTagRequestBody struct{}"
	}

	return strings.Join([]string{"CreateFirewallTagRequestBody", string(data)}, " ")
}
