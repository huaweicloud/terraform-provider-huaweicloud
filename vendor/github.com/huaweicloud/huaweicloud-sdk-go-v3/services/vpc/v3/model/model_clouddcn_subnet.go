package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClouddcnSubnet
type ClouddcnSubnet struct {

	// clouddcn子网ID
	Id string `json:"id"`

	// 功能说明：VPC所属的项目ID
	ProjectId string `json:"project_id"`

	// 功能说明：子网名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：子网描述 取值范围：0-255个字符，不能包含“<”和“>”。
	Description string `json:"description"`

	// 功能说明：子网的网段 取值范围：必须在vpc对应cidr范围内 约束：必须是cidr格式。掩码长度不能大于28
	Cidr string `json:"cidr"`

	// 功能说明：子网的网关 取值范围：子网网段中的IP地址 约束：必须是ip格式
	GatewayIp string `json:"gateway_ip"`

	// clouddcn子网dns服务器地址列表
	DnsNameservers []string `json:"dns_nameservers"`

	// clouddcn子网所在VPC标识
	VpcId string `json:"vpc_id"`

	// 功能说明：资源创建UTC时间 格式：yyyy-MM-ddTHH:mm:ss
	CreatedAt *sdktime.SdkTime `json:"created_at"`

	// 功能说明：资源更新UTC时间 格式：yyyy-MM-ddTHH:mm:ss
	UpdatedAt *sdktime.SdkTime `json:"updated_at"`

	// 功能说明：可用区
	AvailabilityZone string `json:"availability_zone"`

	// 功能说明：对接TMS
	Tags []TagEntity `json:"tags"`

	// 功能说明：对接TPS
	EnterpriseProjectId string `json:"enterprise_project_id"`
}

func (o ClouddcnSubnet) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClouddcnSubnet struct{}"
	}

	return strings.Join([]string{"ClouddcnSubnet", string(data)}, " ")
}
