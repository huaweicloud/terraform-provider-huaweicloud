package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EsLoadBalancerResource struct {

	// 负载均衡器ID。
	Id *string `json:"id,omitempty"`

	// 负载均衡器名称。
	Name *string `json:"name,omitempty"`

	// 是否独享型LB。 - false：共享型。 - true：独享型。
	Guaranteed *string `json:"guaranteed,omitempty"`

	// 资源账单信息 - 空：按需计费。 - 非空：包周期计费。
	BillingInfo *string `json:"billing_info,omitempty"`

	// 描述信息。
	Description *string `json:"description,omitempty"`

	// 负载均衡器所属VPC ID。
	VpcId *string `json:"vpc_id,omitempty"`

	// 负载均衡器的配置状态。
	ProvisioningStatus *string `json:"provisioning_status,omitempty"`

	Listeners *EsListenersResource `json:"listeners,omitempty"`

	// 负载均衡器的IPv4虚拟IP地址。
	VipAddress *string `json:"vip_address,omitempty"`

	Publicips *EsPublicipsResource `json:"publicips,omitempty"`
}

func (o EsLoadBalancerResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EsLoadBalancerResource struct{}"
	}

	return strings.Join([]string{"EsLoadBalancerResource", string(data)}, " ")
}
