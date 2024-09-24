package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutopilotEniNetwork ENI网络配置。
type AutopilotEniNetwork struct {

	// ENI所在子网的IPv4子网ID列表。获取方法如下：  - 方法1：登录虚拟私有云服务的控制台界面，单击VPC下的子网，进入子网详情页面，查找IPv4子网ID。 - 方法2：通过虚拟私有云服务的查询子网列表接口查询。   [链接请参见[查询子网列表](https://support.huaweicloud.com/api-vpc/vpc_subnet01_0003.html)](tag:hws)   [链接请参见[查询子网列表](https://support.huaweicloud.com/intl/zh-cn/api-vpc/vpc_subnet01_0003.html)](tag:hws_hk)
	Subnets []AutopilotNetworkSubnet `json:"subnets"`
}

func (o AutopilotEniNetwork) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotEniNetwork struct{}"
	}

	return strings.Join([]string{"AutopilotEniNetwork", string(data)}, " ")
}
