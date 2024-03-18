package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HostNetwork Node network parameters.
type HostNetwork struct {

	// 用于创建控制节点的VPC的ID。  获取方法如下： - 方法1：登录虚拟私有云服务的控制台界面，在虚拟私有云的详情页面查找VPC ID。 - 方法2：通过虚拟私有云服务的API接口查询。   [链接请参见[查询VPC列表](https://support.huaweicloud.com/api-vpc/vpc_api01_0003.html)](tag:hws)   [链接请参见[查询VPC列表](https://support.huaweicloud.com/intl/zh-cn/api-vpc/vpc_api01_0003.html)](tag:hws_hk)
	Vpc string `json:"vpc"`

	// 用于创建控制节点的subnet的网络ID。获取方法如下：  - 方法1：登录虚拟私有云服务的控制台界面，单击VPC下的子网，进入子网详情页面，查找网络ID。 - 方法2：通过虚拟私有云服务的查询子网列表接口查询。   [链接请参见[查询子网列表](https://support.huaweicloud.com/api-vpc/vpc_subnet01_0003.html)](tag:hws)   [链接请参见[查询子网列表](https://support.huaweicloud.com/intl/zh-cn/api-vpc/vpc_subnet01_0003.html)](tag:hws_hk)
	Subnet string `json:"subnet"`

	// 集群默认的Node节点安全组ID，不指定该字段系统将自动为用户创建默认Node节点安全组，指定该字段时集群将绑定指定的安全组。Node节点安全组需要放通部分端口来保证正常通信。[详细设置请参考[集群安全组规则配置](https://support.huaweicloud.com/cce_faq/cce_faq_00265.html)。](tag:hws)[详细设置请参考[集群安全组规则配置](https://support.huaweicloud.com/intl/zh-cn/cce_faq/cce_faq_00265.html)。](tag:hws_hk)
	SecurityGroup *string `json:"SecurityGroup,omitempty"`
}

func (o HostNetwork) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostNetwork struct{}"
	}

	return strings.Join([]string{"HostNetwork", string(data)}, " ")
}
