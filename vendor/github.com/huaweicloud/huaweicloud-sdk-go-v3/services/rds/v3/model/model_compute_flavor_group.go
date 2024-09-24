package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ComputeFlavorGroup 查询数据库可变更规格接口中，规格组类型
type ComputeFlavorGroup struct {

	// 性能规格，包含以下状态：  normal：通用增强型。 normal2：通用增强Ⅱ型。 armFlavors：鲲鹏通用增强型。 dedicicateNormal（dedicatedNormalLocalssd）：x86独享型。 armLocalssd：鲲鹏通用型。 normalLocalssd：x86通用型。 general：通用型。 dedicated 对于PostgreSQL引擎：独享型
	GroupType string `json:"group_type"`

	ComputeFlavors *ComputeFlavor `json:"compute_flavors"`
}

func (o ComputeFlavorGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ComputeFlavorGroup struct{}"
	}

	return strings.Join([]string{"ComputeFlavorGroup", string(data)}, " ")
}
