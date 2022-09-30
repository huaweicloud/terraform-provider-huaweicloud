package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShrinkNodeReq struct {

	// 需要缩容的节点数量。  - 对节点缩容后，请确保每个节点类型在每个AZ中的数量至少为1。  - 关于跨AZ的集群，在不同AZ中同类型节点个数的差值要小于等于1。  - 关于没有Master节点的集群，每次缩容的数据节点个数(包含冷数据节点和其他类型节点)要小于当前数据节点总数的一半，缩容后的数据节点个数要大于索引的最大副本个数。  - 关于有Master节点的集群，每次缩容的Master节点个数要小于当前Master节点总数的一半，缩容后的Master节点个数必须是奇数且不小于3。
	ReducedNodeNum int32 `json:"reducedNodeNum"`

	// 指定节点类型。 - ess：数据节点。 - ess-cold：冷数据节点。 - ess-client：Client节点。 - ess-master：Master节点。
	Type string `json:"type"`
}

func (o ShrinkNodeReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShrinkNodeReq struct{}"
	}

	return strings.Join([]string{"ShrinkNodeReq", string(data)}, " ")
}
