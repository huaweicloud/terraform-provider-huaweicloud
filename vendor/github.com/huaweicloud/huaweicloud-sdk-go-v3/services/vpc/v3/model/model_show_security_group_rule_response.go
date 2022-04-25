package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowSecurityGroupRuleResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	SecurityGroupRule *SecurityGroupRule `json:"security_group_rule,omitempty"`
	HttpStatusCode    int                `json:"-"`
}

func (o ShowSecurityGroupRuleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowSecurityGroupRuleResponse struct{}"
	}

	return strings.Join([]string{"ShowSecurityGroupRuleResponse", string(data)}, " ")
}
