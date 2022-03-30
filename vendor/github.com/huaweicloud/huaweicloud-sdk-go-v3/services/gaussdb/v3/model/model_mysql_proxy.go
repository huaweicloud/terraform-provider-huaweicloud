package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlProxy struct {
	// Proxy实例id。

	PoolId *string `json:"pool_id,omitempty"`
	// Proxy实例开启状态。  取值范围：closed、open、frozen、opening、closing、enlarging、freezing和unfreezin。

	Status *string `json:"status,omitempty"`
	// Proxy读写分离地址。

	Address *string `json:"address,omitempty"`
	// Proxy端口信息。

	Port *int32 `json:"port,omitempty"`
	// Proxy实例状态。 取值范围：abnormal、normal、creating和deleted。

	PoolStatus *string `json:"pool_status,omitempty"`
	// 延时阈值，单位：秒。

	DelayThresholdInSeconds *int32 `json:"delay_threshold_in_seconds,omitempty"`
	// Elb模式的虚拟ip信息。

	ElbVip *string `json:"elb_vip,omitempty"`
	// 弹性公网IP信息。

	Eip *string `json:"eip,omitempty"`
	// Proxy实例规格的CPU数量。

	Vcpus *string `json:"vcpus,omitempty"`
	// Proxy实例规格的内存数量。

	Ram *string `json:"ram,omitempty"`
	// Proxy节点个数。

	NodeNum *int32 `json:"node_num,omitempty"`
	// Proxy主备模式，取值范围：Cluster。

	Mode *string `json:"mode,omitempty"`
	// Proxy节点信息。

	Nodes *[]MysqlProxyNodes `json:"nodes,omitempty"`
	// Proxy规格信息。

	FlavorRef *string `json:"flavor_ref,omitempty"`
}

func (o MysqlProxy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlProxy struct{}"
	}

	return strings.Join([]string{"MysqlProxy", string(data)}, " ")
}
