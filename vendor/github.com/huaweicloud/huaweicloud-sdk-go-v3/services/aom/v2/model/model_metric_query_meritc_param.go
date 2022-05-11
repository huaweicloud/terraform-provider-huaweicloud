package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 查询参数集
type MetricQueryMeritcParam struct {

	// 指标维度列表。 取值范围： 数组不能为空，同时数组中任何一个dimension对象name和value属性的值也不能为空。
	Dimensions []Dimension `json:"dimensions"`

	// 指标名称。名称长度取值范围为1~255个字符。 取值范围： AOM提供的基础指标， cpuUsage、cpuCoreUsed等 cpuUage：cpu使用率； cpuCoreUsed：cpu内核占用； 用户上报的自定义指标名称。
	MetricName string `json:"metricName"`

	// 指标命名空间。 取值范围： PAAS.CONTAINER：组件指标、实例指标、进程指标和容器指标的命名空间， PAAS.NODE： 主机指标、网络指标、磁盘指标和文件系统指标的命名空间， PAAS.SLA：SLA指标的命名空间， PAAS.AGGR：集群指标的命名空间， CUSTOMMETRICS：默认的自定义指标的命名空间。
	Namespace string `json:"namespace"`
}

func (o MetricQueryMeritcParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MetricQueryMeritcParam struct{}"
	}

	return strings.Join([]string{"MetricQueryMeritcParam", string(data)}, " ")
}
