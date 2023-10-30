package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShrinkNodesReq struct {

	// 是否迁移数据。 - \"true\"：迁移数据。 - \"false\"：不迁移数据。
	MigrateData *string `json:"migrate_data,omitempty"`

	// 需要缩容的节点ID。  通过[查询集群详情](ShowClusterDetail.xml)获取instances中的id属性。
	ShrinkNodes []string `json:"shrinkNodes"`
}

func (o ShrinkNodesReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShrinkNodesReq struct{}"
	}

	return strings.Join([]string{"ShrinkNodesReq", string(data)}, " ")
}
