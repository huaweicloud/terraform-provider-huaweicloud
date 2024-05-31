package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClouddcnSubnetOption
type CreateClouddcnSubnetOption struct {

	// 功能说明：clouddcn子网名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：clouddcn子网描述 取值范围：0-255个字符，不能包含“<”和“>”。
	Description *string `json:"description,omitempty"`

	// 功能说明：clouddcn子网的网段 取值范围：必须在vpc对应cidr范围内，不能和同vpc下其他普通子网的网段冲突 约束：必须是cidr格式。掩码长度不能大于28
	Cidr string `json:"cidr"`

	// clouddcn子网所在VPC标识
	VpcId string `json:"vpc_id"`

	// 功能说明：clouddcn子网的网关 取值范围：clouddcn子网网段中的IP地址 约束：必须是ip格式
	GatewayIp string `json:"gateway_ip"`

	// 功能说明：clouddcn子网dns服务器地址的集合；如果想使用两个以上dns服务器，请使用该字段 约束：是子网dns服务器地址1跟子网dns服务器地址2的合集的父集，不支持IPv6地址。 默认值：不填时为空，无法使用云内网DNS功能 [内网DNS地址请参见](https://support.huaweicloud.com/dns_faq/dns_faq_002.html) [通过API获取请参见](https://support.huaweicloud.com/api-dns/dns_api_69001.html)
	DnsNameservers *[]string `json:"dns_nameservers,omitempty"`

	// 功能说明：可用区
	AvailabilityZone *string `json:"availability_zone,omitempty"`

	// 功能说明：对接TMS
	Tags *[]TagEntity `json:"tags,omitempty"`
}

func (o CreateClouddcnSubnetOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClouddcnSubnetOption struct{}"
	}

	return strings.Join([]string{"CreateClouddcnSubnetOption", string(data)}, " ")
}
