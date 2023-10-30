package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateFirewallRequestBody This is a auto create Body Object
type CreateFirewallRequestBody struct {
	Firewall *CreateFirewallOption `json:"firewall"`

	// 功能说明：是否只预检此次请求 取值范围： -true：发送检查请求，不会创建ACL。检查项包括是否填写了必需参数、请求格式、业务限制。如果检查不通过，则返回对应错误。如果检查通过，则返回响应码202。 -false（默认值）：发送正常请求，并直接创建ACL。
	DryRun *bool `json:"dry_run,omitempty"`
}

func (o CreateFirewallRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateFirewallRequestBody struct{}"
	}

	return strings.Join([]string{"CreateFirewallRequestBody", string(data)}, " ")
}
