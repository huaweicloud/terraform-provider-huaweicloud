package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FirewallUpdateRuleItemOption
type FirewallUpdateRuleItemOption struct {

	// 功能说明：ACL规则唯一标识 取值范围：合法UUID的字符串
	Id string `json:"id"`

	// 功能说明：ACL规则名称 取值范围：0-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name *string `json:"name,omitempty"`

	// 功能说明：ACL规则描述信息 取值范围：0-255个字符 约束：不能包含“<”和“>”。
	Description *string `json:"description,omitempty"`

	// 功能说明：ACL规则对流量执行的操作放通或拒绝 取值范围：allow放通；deny拒绝
	Action *string `json:"action,omitempty"`

	// 功能说明：ACL规则协议 取值范围：支持tcp,udp,icmp,icmpv6或者协议号（0-255），any表示全部协议
	Protocol *string `json:"protocol,omitempty"`

	// 功能说明：ACL规则的ip版本 取值范围：4, 表示ipv4；6, 表示ipv6
	IpVersion *int32 `json:"ip_version,omitempty"`

	// 功能说明：ACL规则源IP地址或者CIDR 约束：source_ip_address和source_address_group_id不能同时设置
	SourceIpAddress *string `json:"source_ip_address,omitempty"`

	// 功能说明：ACL规则目的IP地址或者CIDR 约束：destination_ip_address和destination_address_group_id不能同时设置
	DestinationIpAddress *string `json:"destination_ip_address,omitempty"`

	// 功能说明：ACL规则的源端口 取值范围：支持端口号，一段端口范围，多个以逗号分隔 约束：支持的端口组的数量默认为20
	SourcePort *string `json:"source_port,omitempty"`

	// 功能说明：ACL规则的目的端口 取值范围：支持端口号，一段端口范围，多个以逗号分隔 约束：支持的端口组的数量默认为20
	DestinationPort *string `json:"destination_port,omitempty"`

	// 功能说明：ACL规则的源地址组ID 约束：source_ip_address和source_address_group_id不能同时设置
	SourceAddressGroupId *string `json:"source_address_group_id,omitempty"`

	// 功能说明：ACL规则的目的地址组ID 约束：destination_ip_address和destination_address_group_id不能同时设置
	DestinationAddressGroupId *string `json:"destination_address_group_id,omitempty"`

	// 功能说明：是否启用ACL规则。 取值范围：true，表示启用；false，表示停用。 默认值：true。
	Enabled *bool `json:"enabled,omitempty"`
}

func (o FirewallUpdateRuleItemOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FirewallUpdateRuleItemOption struct{}"
	}

	return strings.Join([]string{"FirewallUpdateRuleItemOption", string(data)}, " ")
}
