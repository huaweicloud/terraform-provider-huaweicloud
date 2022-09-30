package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowCeshierarchyRespDimensions struct {

	// 监控维度名称。
	Name *string `json:"name,omitempty"`

	// 监控指标名称。
	Metrics *[]string `json:"metrics,omitempty"`

	// 监控查询使用的key。
	KeyName *[]string `json:"key_name,omitempty"`

	// 监控维度路由。
	DimRouter *[]string `json:"dim_router,omitempty"`

	// 子维度列表。
	Children *[]ShowCeshierarchyRespChildren `json:"children,omitempty"`
}

func (o ShowCeshierarchyRespDimensions) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCeshierarchyRespDimensions struct{}"
	}

	return strings.Join([]string{"ShowCeshierarchyRespDimensions", string(data)}, " ")
}
