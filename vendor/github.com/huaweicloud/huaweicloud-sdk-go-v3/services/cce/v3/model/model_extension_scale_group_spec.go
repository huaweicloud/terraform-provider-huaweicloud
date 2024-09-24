package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExtensionScaleGroupSpec 扩展伸缩组配置，承载区别于默认伸缩组的差异化配置
type ExtensionScaleGroupSpec struct {

	// 节点规格
	Flavor *string `json:"flavor,omitempty"`

	// 节点可用区，未指定或者为空则以默认伸缩组中配置为准
	Az *string `json:"az,omitempty"`

	CapacityReservationSpecification *CapacityReservationSpecification `json:"capacityReservationSpecification,omitempty"`

	Autoscaling *ScaleGroupAutoscaling `json:"autoscaling,omitempty"`
}

func (o ExtensionScaleGroupSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExtensionScaleGroupSpec struct{}"
	}

	return strings.Join([]string{"ExtensionScaleGroupSpec", string(data)}, " ")
}
