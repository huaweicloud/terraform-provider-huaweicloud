package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateSecurityGroupRulesResponse Response Object
type BatchCreateSecurityGroupRulesResponse struct {

	// 批量创建安全组规则的响应体
	SecurityGroupRules *[]SecurityGroupRule `json:"security_group_rules,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o BatchCreateSecurityGroupRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateSecurityGroupRulesResponse struct{}"
	}

	return strings.Join([]string{"BatchCreateSecurityGroupRulesResponse", string(data)}, " ")
}
