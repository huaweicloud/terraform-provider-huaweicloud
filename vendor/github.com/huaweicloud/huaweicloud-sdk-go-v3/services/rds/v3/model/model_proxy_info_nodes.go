package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProxyInfoNodes struct {

	// 数据库代理节点ID。
	Id *string `json:"id,omitempty"`

	// 数据库代理节点状态。  取值范围： NORMAL: 表示节点正常。 ABNORMAL: 表示节点节点状态异常。 CREATING: 表示节点正在创建中。 CREATEFAIL: 表示节点创建失败。
	Status *string `json:"status,omitempty"`

	// 数据库代理节点角色：  master：主节点。  slave：备节点。
	Role *string `json:"role,omitempty"`

	// 数据库代理节点所在可用区。
	AzCode *string `json:"az_code,omitempty"`

	// 数据库代理节点是否被冻结。  取值范围：  0：未冻结。  1：冻结。
	FrozenFlag *int32 `json:"frozen_flag,omitempty"`
}

func (o ProxyInfoNodes) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProxyInfoNodes struct{}"
	}

	return strings.Join([]string{"ProxyInfoNodes", string(data)}, " ")
}
