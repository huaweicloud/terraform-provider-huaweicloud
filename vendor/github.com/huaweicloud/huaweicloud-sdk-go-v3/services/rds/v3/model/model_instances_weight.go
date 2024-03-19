package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InstancesWeight 数据库节点的读权重设置。  在proxy_mode为readonly时，只能为只读节点选择权重。
type InstancesWeight struct {

	// 数据库实例ID。
	InstanceId string `json:"instance_id"`

	// 数据库代理读权重。
	Weight int32 `json:"weight"`
}

func (o InstancesWeight) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstancesWeight struct{}"
	}

	return strings.Join([]string{"InstancesWeight", string(data)}, " ")
}
