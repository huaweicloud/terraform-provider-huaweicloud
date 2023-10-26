package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateFirewallOption
type UpdateFirewallOption struct {

	// 功能说明：ACL名称 取值范围：0-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name *string `json:"name,omitempty"`

	// 功能说明：ACL描述信息 取值范围：0-255个字符 约束：不能包含“<”和“>”
	Description *string `json:"description,omitempty"`

	// 功能说明：ACL是否开启 取值范围：true表示ACL开启；false表示ACL关闭
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

func (o UpdateFirewallOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFirewallOption struct{}"
	}

	return strings.Join([]string{"UpdateFirewallOption", string(data)}, " ")
}
