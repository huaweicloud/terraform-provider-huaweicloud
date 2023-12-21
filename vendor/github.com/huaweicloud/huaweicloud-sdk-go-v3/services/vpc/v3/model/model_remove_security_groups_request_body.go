package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveSecurityGroupsRequestBody This is a auto create Body Object
type RemoveSecurityGroupsRequestBody struct {
	Port *RemoveSecurityGroupOption `json:"port"`
}

func (o RemoveSecurityGroupsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveSecurityGroupsRequestBody struct{}"
	}

	return strings.Join([]string{"RemoveSecurityGroupsRequestBody", string(data)}, " ")
}
