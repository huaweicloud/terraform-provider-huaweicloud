package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListFlavorsResizeResponse Response Object
type ListFlavorsResizeResponse struct {

	// 规格组列表  normal：通用增强型。 normal2：通用增强Ⅱ型。 armFlavors：鲲鹏通用增强型。 dedicicateNormal（dedicatedNormalLocalssd）：x86独享型。 armLocalssd：鲲鹏通用型。 normalLocalssd：x86通用型。 general：通用型。 dedicated 对于PostgreSQL引擎：独享型
	FlavorGroups   *[]ComputeFlavorGroup `json:"flavor_groups,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o ListFlavorsResizeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFlavorsResizeResponse struct{}"
	}

	return strings.Join([]string{"ListFlavorsResizeResponse", string(data)}, " ")
}
