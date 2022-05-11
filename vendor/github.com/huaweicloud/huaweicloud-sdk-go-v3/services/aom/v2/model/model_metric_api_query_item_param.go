package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 指标查询参数详情。
type MetricApiQueryItemParam struct {

	// 资源编号,格式为resType_resId。其中resType部分的枚举值为：host, application, instance, container, process, network, storage, volume。当URI参数中的type取值为“inventory”时，通过该参数查询关联的指标，不再使用metricItems数组中的信息。
	InventoryId *string `json:"inventoryId,omitempty"`

	// 当URI参数中的type取值不为“inventory”时，就通过该数组传递的参数信息进行指标查询。
	MetricItems *[]QueryMetricItemOptionParam `json:"metricItems,omitempty"`
}

func (o MetricApiQueryItemParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MetricApiQueryItemParam struct{}"
	}

	return strings.Join([]string{"MetricApiQueryItemParam", string(data)}, " ")
}
