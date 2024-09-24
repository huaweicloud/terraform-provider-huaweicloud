package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutopilotNetworkSubnet ENI网络配置，创建集群指定subnets字段使用时必填。
type AutopilotNetworkSubnet struct {

	// 用于创建控制节点的subnet的IPv4子网ID(暂不支持IPv6)。获取方法如下：  - 方法1：登录虚拟私有云服务的控制台界面，单击VPC下的子网，进入子网详情页面，查找IPv4子网ID。 - 方法2：通过虚拟私有云服务的查询子网列表接口查询。   [链接请参见[查询子网列表](https://support.huaweicloud.com/api-vpc/vpc_subnet01_0003.html)](tag:hws)   [链接请参见[查询子网列表](https://support.huaweicloud.com/intl/zh-cn/api-vpc/vpc_subnet01_0003.html)](tag:hws_hk)
	SubnetID string `json:"subnetID"`
}

func (o AutopilotNetworkSubnet) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotNetworkSubnet struct{}"
	}

	return strings.Join([]string{"AutopilotNetworkSubnet", string(data)}, " ")
}
