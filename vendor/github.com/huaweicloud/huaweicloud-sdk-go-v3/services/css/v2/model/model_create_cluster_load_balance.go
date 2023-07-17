package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterLoadBalance 终端节点服务信息。
type CreateClusterLoadBalance struct {

	// 是否开启内网域名。 - true： 开启内网域名。 - false： 关闭内网域名。
	EndpointWithDnsName bool `json:"endpointWithDnsName"`

	// 访问控制。
	VpcPermissions *[]string `json:"vpcPermissions,omitempty"`
}

func (o CreateClusterLoadBalance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterLoadBalance struct{}"
	}

	return strings.Join([]string{"CreateClusterLoadBalance", string(data)}, " ")
}
