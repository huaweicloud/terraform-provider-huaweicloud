package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlProxyNodes struct {
	// Proxy节点id。

	Id *string `json:"id,omitempty"`
	// Proxy节点状态。 取值范围：normal、abnormal、creating和deleted。

	Status *string `json:"status,omitempty"`
	// Proxy节点名称。

	Name *string `json:"name,omitempty"`
	// Proxy节点角色：master和slave。

	Role *string `json:"role,omitempty"`
	// 可用区。

	AzCode *string `json:"az_code,omitempty"`
	// Proxy节点是否被冻结：0-未冻结；1-冻结；2-冻结删除。

	FrozenFlag *int32 `json:"frozen_flag,omitempty"`
}

func (o MysqlProxyNodes) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlProxyNodes struct{}"
	}

	return strings.Join([]string{"MysqlProxyNodes", string(data)}, " ")
}
