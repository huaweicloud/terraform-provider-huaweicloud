package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RollingRestartReq struct {

	// 操作角色。type的值只能为role。
	Type string `json:"type"`

	// 实例类型（选择实例类型时至少需要一个数据节点），多个类型使用逗号隔开。例如:  - ess-master对应Master节点。 - ess-client对应Client节点。 - ess-cold对应冷数据节点。 - ess对应数据节点。 - all对应所有节点。
	Value string `json:"value"`
}

func (o RollingRestartReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RollingRestartReq struct{}"
	}

	return strings.Join([]string{"RollingRestartReq", string(data)}, " ")
}
