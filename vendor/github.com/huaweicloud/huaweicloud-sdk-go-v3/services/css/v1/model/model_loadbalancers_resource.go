package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type LoadbalancersResource struct {

	// 负载均衡器ID。
	Id *string `json:"id,omitempty"`

	// 负载均衡器名称。
	Name *string `json:"name,omitempty"`

	// 7层协议Id。
	L7FlavorId *string `json:"l7_flavor_id,omitempty"`

	// 是否开启跨VPC后端。
	IpTargetEnable *bool `json:"ip_target_enable,omitempty"`
}

func (o LoadbalancersResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LoadbalancersResource struct{}"
	}

	return strings.Join([]string{"LoadbalancersResource", string(data)}, " ")
}
