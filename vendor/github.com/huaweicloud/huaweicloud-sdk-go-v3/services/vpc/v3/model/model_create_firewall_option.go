package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateFirewallOption
type CreateFirewallOption struct {

	// 功能说明：ACL名称 取值范围：0-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：ACL描述信息 取值范围：0-255个字符 约束：不能包含“<”和“>”。
	Description *string `json:"description,omitempty"`

	// 功能说明：ACL企业项目ID。 取值范围：最大长度36字节，带“-”连字符的UUID格式，或者是字符串“0”。“0”表示默认企业项目。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 功能描述：ACL资源标签
	Tags *[]ResourceTag `json:"tags,omitempty"`

	// 功能说明：ACL是否开启，默认值true 取值范围：true表示ACL开启；false表示ACL关闭
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

func (o CreateFirewallOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateFirewallOption struct{}"
	}

	return strings.Join([]string{"CreateFirewallOption", string(data)}, " ")
}
