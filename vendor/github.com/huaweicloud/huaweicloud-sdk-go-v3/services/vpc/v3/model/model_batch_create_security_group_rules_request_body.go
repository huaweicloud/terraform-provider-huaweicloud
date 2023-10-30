package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateSecurityGroupRulesRequestBody This is a auto create Body Object
type BatchCreateSecurityGroupRulesRequestBody struct {

	// 待创建的安全组规则列表
	SecurityGroupRules []BatchCreateSecurityGroupRulesOption `json:"security_group_rules"`

	// 创建时是否忽略重复的安全组规则 默认为false
	IgnoreDuplicate *bool `json:"ignore_duplicate,omitempty"`

	// 功能说明：是否只预检此次请求 取值范围： -true：发送检查请求，不会创建安全组规则。检查项包括是否填写了必需参数、请求格式、业务限制。如果检查不通过，则返回对应错误。如果检查通过，则返回响应码202。 -false（默认值）：发送正常请求，并直接创建安全组规则。
	DryRun *bool `json:"dry_run,omitempty"`
}

func (o BatchCreateSecurityGroupRulesRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateSecurityGroupRulesRequestBody struct{}"
	}

	return strings.Join([]string{"BatchCreateSecurityGroupRulesRequestBody", string(data)}, " ")
}
