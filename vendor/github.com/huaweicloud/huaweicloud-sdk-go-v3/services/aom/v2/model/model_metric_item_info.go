package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 指标详细信息。
type MetricItemInfo struct {

	// 指标维度列表。 维度最多允许50个，单个维度为json对象，结构说明如下 dimension.name：长度最短为1，最大为32。 dimension.value：长度最短为1，最大为64。
	Dimensions []Dimension2 `json:"dimensions"`

	// 指标命名空间。 namespace中不允许存在\":\"符号，取值范围 格式为service.item；service和item必须是字符串，必须以字母开头，只能包含0-9/a-z/A-Z/_，总长度最短为3，最大为32，service不能为“PAAS”。
	Namespace string `json:"namespace"`
}

func (o MetricItemInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MetricItemInfo struct{}"
	}

	return strings.Join([]string{"MetricItemInfo", string(data)}, " ")
}
