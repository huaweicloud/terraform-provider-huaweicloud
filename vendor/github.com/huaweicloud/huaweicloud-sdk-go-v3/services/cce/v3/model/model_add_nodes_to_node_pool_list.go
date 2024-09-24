package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddNodesToNodePoolList 自定义节点池纳管节点参数，纳管过程将清理节点上系统盘、数据盘数据，并作为新节点接入Kuberntes集群，请提前备份迁移关键数据。
type AddNodesToNodePoolList struct {

	// API版本，固定值“v3”。
	ApiVersion string `json:"apiVersion"`

	// API类型，固定值“List”。
	Kind string `json:"kind"`

	// 纳管节点列表，当前最多支持同时纳管200个节点。
	NodeList []AddNodesToNodePool `json:"nodeList"`
}

func (o AddNodesToNodePoolList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddNodesToNodePoolList struct{}"
	}

	return strings.Join([]string{"AddNodesToNodePoolList", string(data)}, " ")
}
