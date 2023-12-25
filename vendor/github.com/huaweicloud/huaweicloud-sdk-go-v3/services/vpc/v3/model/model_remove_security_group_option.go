package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RemoveSecurityGroupOption
type RemoveSecurityGroupOption struct {

	// 功能说明：安全组的ID列表；例如：\"security_groups\": [\"a0608cbf-d047-4f54-8b28-cd7b59853fff\"]
	SecurityGroups []string `json:"security_groups"`
}

func (o RemoveSecurityGroupOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveSecurityGroupOption struct{}"
	}

	return strings.Join([]string{"RemoveSecurityGroupOption", string(data)}, " ")
}
