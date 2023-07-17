package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RestartClusterReq struct {

	// 操作角色。参数范围： - node - role
	Type string `json:"type"`

	// 操作参数。参数说明：  - 当操作角色为node时，value为节点ID,通过[查询集群详情](ShowClusterDetail.xml)获取instances中的id属性。  - 当操作角色为role时，value为节点类型(ess、ess-master、ess-client、ess-cold)的多种不同组合。
	Value string `json:"value"`
}

func (o RestartClusterReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestartClusterReq struct{}"
	}

	return strings.Join([]string{"RestartClusterReq", string(data)}, " ")
}
